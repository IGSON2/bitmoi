-- name: CreateVerifyEmail :execresult
INSERT INTO verify_emails (
    user_id,
    secret_code,
    created_at,
    expired_at
) VALUES (
    ?, ?, ?, ?
);

-- name: UpdateVerifyEmail :execresult
UPDATE verify_emails
SET
    is_used = TRUE
WHERE
    id = ?
    AND secret_code = ?
    AND is_used = FALSE
    AND expired_at > now();

-- name: GetVerifyEmails :one
SELECT * FROM verify_emails
WHERE id = ? AND secret_code= ?;