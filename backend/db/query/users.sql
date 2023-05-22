-- name: CreateUser :execresult
INSERT INTO users (
    user_id,
    uid,
    fullname,
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

-- name: UpdatePhotoURL :execresult
UPDATE users SET photo_url = ?
WHERE user_id = ?;