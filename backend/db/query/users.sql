-- name: CreateUser :execresult
INSERT INTO users (
    user_id,
    oauth_uid,
    nickname,
    hashed_password,
    email,
    photo_url
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = ?;

-- name: GetUserByNickName :one
SELECT * FROM users
WHERE nickname = ?;

-- name: GetRandomUser :one
SELECT * FROM users
ORDER BY RAND()
LIMIT 1;

-- name: GetLastUser :one
SELECT * FROM users
ORDER BY created_at DESC
LIMIT 1;

-- name: UpdateUserPhotoURL :execresult
UPDATE users SET photo_url = ?
WHERE user_id = ?;

-- name: UpdateUserEmailVerified :execresult
UPDATE users SET is_email_verified = ?
WHERE user_id = ?;