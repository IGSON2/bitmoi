package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/mail"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	taskTest      = "task:test"
)

type TestRedisTaskProcessor struct {
	server         *asynq.Server
	store          db.Store
	mailer         mail.EmailSender
	accessDuration time.Duration
}

func NewTestRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender, accessDuration time.Duration) *TestRedisTaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: worker.NewLogger(),
	})

	return &TestRedisTaskProcessor{
		server:         server,
		store:          store,
		mailer:         mailer,
		accessDuration: accessDuration,
	}
}

func (processor *TestRedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(taskTest, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

func (processor *TestRedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload worker.PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("faild to unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	secretCode := utilities.MakeRanString(32)

	r, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		UserID:     user.UserID,
		Email:      user.Email,
		SecretCode: secretCode,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(processor.accessDuration),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("cannot get id for last inserted row: %w", err)
	}

	verifyEmail, err := processor.store.GetVerifyEmails(ctx, db.GetVerifyEmailsParams{
		ID:         id,
		SecretCode: secretCode,
	})
	if err != nil {
		return fmt.Errorf("cannot get verifyEmail by specified data: %w", err)
	}

	verifyUrl := fmt.Sprintf("http://localhost:5001/user/verifyEmail?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)

	client := &http.Client{}
	req, err := http.NewRequest("GET", verifyUrl, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	vr := new(VerifyEmailResponse)
	json.Unmarshal(b, vr)

	if !vr.IsVerified {
		return fmt.Errorf("not verified")
	}

	return nil
}

type RedisTestTaskDistributor struct {
	client *asynq.Client
}

func (distributor *RedisTestTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *worker.PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload :%w", err)
	}

	task := asynq.NewTask(taskTest, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task :%w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) worker.TaskDistributor {
	client := asynq.NewClient(redisOpt)

	return &RedisTestTaskDistributor{
		client: client,
	}
}
