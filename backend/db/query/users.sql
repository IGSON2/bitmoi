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

-- name: GetUserByMetamaskAddress :one
SELECT * FROM users
WHERE metamask_address = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

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

-- name: UpdateUserMetamaskAddress :execresult
UPDATE users 
SET 
    metamask_address = ?,
    address_changed_at = ?
WHERE user_id = ?;