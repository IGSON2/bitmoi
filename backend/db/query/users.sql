-- name: CreateUser :execresult
INSERT INTO users (
    user_id,
    oauth_uid,
    full_name,
    hashed_password,
    email,
    photo_url
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ?
LIMIT 1;

-- name: GetRandomUser :one
SELECT * FROM users
ORDER BY RAND()
LIMIT 1;

-- name: GetLastUser :one
SELECT * FROM users
ORDER BY created_at DESC
LIMIT 1;

-- name: GetUserMetamaskAddress :one
SELECT metamask_address FROM users
WHERE user_id = ?;

-- name: UpdateUserPhotoURL :execresult
UPDATE users SET photo_url = ?
WHERE user_id = ?;

-- name: UpdateUserMetamaskAddress :execresult
UPDATE users SET metamask_address = ?
WHERE user_id = ?;

-- name: UpdateUserIsEmailVerified :execresult
UPDATE users SET is_email_verified = ?
WHERE user_id = ?;