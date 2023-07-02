// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: verify_emails.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :execresult
INSERT INTO verify_emails (
    user_id,
    email,
    secret_code,
    created_at,
    expired_at
) VALUES (
    ?, ?, ?, ?, ?
)
`

type CreateVerifyEmailParams struct {
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createVerifyEmail,
		arg.UserID,
		arg.Email,
		arg.SecretCode,
		arg.CreatedAt,
		arg.ExpiredAt,
	)
}

const getVerifyEmails = `-- name: GetVerifyEmails :one
SELECT id, user_id, email, secret_code, is_used, created_at, expired_at FROM verify_emails
WHERE id = ? AND secret_code= ?
`

type GetVerifyEmailsParams struct {
	ID         int64  `json:"id"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) GetVerifyEmails(ctx context.Context, arg GetVerifyEmailsParams) (VerifyEmail, error) {
	row := q.db.QueryRowContext(ctx, getVerifyEmails, arg.ID, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateVerifyEmail = `-- name: UpdateVerifyEmail :execresult
UPDATE verify_emails
SET
    is_used = TRUE
WHERE
    id = ?
    AND secret_code = ?
    AND is_used = FALSE
    AND expired_at > now()
`

type UpdateVerifyEmailParams struct {
	ID         int64  `json:"id"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateVerifyEmail, arg.ID, arg.SecretCode)
}
