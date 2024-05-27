package worker

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/mail"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server         *asynq.Server
	store          db.Store
	mailer         mail.EmailSender
	accessDuration time.Duration
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender, accessDuration time.Duration) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: logger,
	})

	return &RedisTaskProcessor{
		server:         server,
		store:          store,
		mailer:         mailer,
		accessDuration: accessDuration,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
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

	subject := "[BITMOI] 인증 메일 안내"
	verifyUrl := fmt.Sprintf("https://api.bitmoi.co.kr/basic/user/verifyEmail?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := utilities.GenerateEmailMessage(user.UserID, verifyUrl)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")

	return nil
}
