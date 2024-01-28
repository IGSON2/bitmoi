-- name: CreateUser :execresult
INSERT INTO users (
    user_id,
    oauth_uid,
    nickname,
    hashed_password,
    email,
    photo_url,
    prac_balance,
    comp_balance,
    recommender_code
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
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

-- name: GetLastUserID :one
SELECT id FROM users
ORDER BY id DESC
LIMIT 1;

-- name: GetUserLastAccessedAt :one
SELECT last_accessed_at FROM users
WHERE user_id = ?;

-- name: GetUserPracBalance :one
SELECT prac_balance FROM users
WHERE user_id = ?;

-- name: UpdateUserPhotoURL :execresult
UPDATE users SET photo_url = ?
WHERE user_id = ?;

-- name: UpdateUserMetamaskAddress :execresult
UPDATE users 
SET 
    metamask_address = ?,
    address_changed_at = ?
WHERE user_id = ?;

-- name: UpdateUserPracBalance :execresult
UPDATE users SET prac_balance = ?
WHERE user_id = ?;

-- name: UpdateUserCompBalance :execresult
UPDATE users SET comp_balance = ?
WHERE user_id = ?;

-- name: UpdateUserLastAccessedAt :execresult
UPDATE users SET last_accessed_at = ?
WHERE user_id = ?;
