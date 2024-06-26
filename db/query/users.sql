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

-- name: GetUsers :many
SELECT * FROM users LIMIT ? OFFSET ?;

-- name: GetUserByNickName :one
SELECT * FROM users
WHERE nickname = ?;

-- name: GetUserByMetamaskAddress :one
SELECT * FROM users
WHERE metamask_address = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: GetUserByRecommenderCode :one
SELECT * FROM users
WHERE recommender_code = ?;

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

-- name: UpdateUserNickname :execresult
UPDATE users 
SET 
    nickname = ?
WHERE user_id = ?;

-- name: AppendUserPracBalance :execresult
UPDATE users SET prac_balance = prac_Balance + ?
WHERE user_id = ?;

-- name: AppendUserCompBalance :execresult
UPDATE users SET comp_balance = comp_balance + ?
WHERE user_id = ?;

-- name: AppendUserWmoiBalance :execresult
UPDATE users SET wmoi_balance = wmoi_balance + ?
WHERE user_id = ?;

-- name: UpdateUserWmoiBalanceByRecom :execresult
UPDATE users SET wmoi_balance = ?
WHERE user_id = ?;

-- name: UpdateUserLastAccessedAt :execresult
UPDATE users SET last_accessed_at = ?
WHERE user_id = ?;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE user_id = ?;