package db

import (
	"context"
	"fmt"
	"time"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (s *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User, err = q.GetUser(ctx, arg.UserID)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})

	return result, err
}

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SqlStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.VerifyEmail, err = q.GetVerifyEmails(ctx, GetVerifyEmailsParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		if time.Now().After(result.VerifyEmail.ExpiredAt) {
			return fmt.Errorf("verification has expired")
		}

		_, err = q.UpdateUserEmailVerified(ctx, UpdateUserEmailVerifiedParams{
			IsEmailVerified: true,
			UserID:          result.VerifyEmail.UserID,
		})
		if err != nil {
			return err
		}

		result.User, err = q.GetUser(ctx, result.VerifyEmail.UserID)
		return err
	})

	return result, err
}
