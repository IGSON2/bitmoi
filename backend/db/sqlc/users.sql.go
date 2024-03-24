// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
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
)
`

type CreateUserParams struct {
	UserID          string         `json:"user_id"`
	OauthUid        sql.NullString `json:"oauth_uid"`
	Nickname        string         `json:"nickname"`
	HashedPassword  sql.NullString `json:"hashed_password"`
	Email           string         `json:"email"`
	PhotoUrl        sql.NullString `json:"photo_url"`
	PracBalance     float64        `json:"prac_balance"`
	CompBalance     float64        `json:"comp_balance"`
	RecommenderCode string         `json:"recommender_code"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser,
		arg.UserID,
		arg.OauthUid,
		arg.Nickname,
		arg.HashedPassword,
		arg.Email,
		arg.PhotoUrl,
		arg.PracBalance,
		arg.CompBalance,
		arg.RecommenderCode,
	)
}

const getLastUserID = `-- name: GetLastUserID :one
SELECT id FROM users
ORDER BY id DESC
LIMIT 1
`

func (q *Queries) GetLastUserID(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLastUserID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getRandomUser = `-- name: GetRandomUser :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
ORDER BY RAND()
LIMIT 1
`

func (q *Queries) GetRandomUser(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getRandomUser)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
WHERE user_id = ?
`

func (q *Queries) GetUser(ctx context.Context, userID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUserByMetamaskAddress = `-- name: GetUserByMetamaskAddress :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
WHERE metamask_address = ?
`

func (q *Queries) GetUserByMetamaskAddress(ctx context.Context, metamaskAddress sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByMetamaskAddress, metamaskAddress)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUserByNickName = `-- name: GetUserByNickName :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
WHERE nickname = ?
`

func (q *Queries) GetUserByNickName(ctx context.Context, nickname string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByNickName, nickname)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUserByRecommenderCode = `-- name: GetUserByRecommenderCode :one
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users
WHERE recommender_code = ?
`

func (q *Queries) GetUserByRecommenderCode(ctx context.Context, recommenderCode string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByRecommenderCode, recommenderCode)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OauthUid,
		&i.Nickname,
		&i.HashedPassword,
		&i.Email,
		&i.MetamaskAddress,
		&i.PhotoUrl,
		&i.PracBalance,
		&i.CompBalance,
		&i.WmoiBalance,
		&i.RecommenderCode,
		&i.CreatedAt,
		&i.LastAccessedAt,
		&i.PasswordChangedAt,
		&i.AddressChangedAt,
	)
	return i, err
}

const getUserLastAccessedAt = `-- name: GetUserLastAccessedAt :one
SELECT last_accessed_at FROM users
WHERE user_id = ?
`

func (q *Queries) GetUserLastAccessedAt(ctx context.Context, userID string) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, getUserLastAccessedAt, userID)
	var last_accessed_at sql.NullTime
	err := row.Scan(&last_accessed_at)
	return last_accessed_at, err
}

const getUserPracBalance = `-- name: GetUserPracBalance :one
SELECT prac_balance FROM users
WHERE user_id = ?
`

func (q *Queries) GetUserPracBalance(ctx context.Context, userID string) (float64, error) {
	row := q.db.QueryRowContext(ctx, getUserPracBalance, userID)
	var prac_balance float64
	err := row.Scan(&prac_balance)
	return prac_balance, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, user_id, oauth_uid, nickname, hashed_password, email, metamask_address, photo_url, prac_balance, comp_balance, wmoi_balance, recommender_code, created_at, last_accessed_at, password_changed_at, address_changed_at FROM users LIMIT ? OFFSET ?
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.OauthUid,
			&i.Nickname,
			&i.HashedPassword,
			&i.Email,
			&i.MetamaskAddress,
			&i.PhotoUrl,
			&i.PracBalance,
			&i.CompBalance,
			&i.WmoiBalance,
			&i.RecommenderCode,
			&i.CreatedAt,
			&i.LastAccessedAt,
			&i.PasswordChangedAt,
			&i.AddressChangedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserCompBalance = `-- name: UpdateUserCompBalance :execresult
UPDATE users SET comp_balance = ?
WHERE user_id = ?
`

type UpdateUserCompBalanceParams struct {
	CompBalance float64 `json:"comp_balance"`
	UserID      string  `json:"user_id"`
}

func (q *Queries) UpdateUserCompBalance(ctx context.Context, arg UpdateUserCompBalanceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserCompBalance, arg.CompBalance, arg.UserID)
}

const updateUserLastAccessedAt = `-- name: UpdateUserLastAccessedAt :execresult
UPDATE users SET last_accessed_at = ?
WHERE user_id = ?
`

type UpdateUserLastAccessedAtParams struct {
	LastAccessedAt sql.NullTime `json:"last_accessed_at"`
	UserID         string       `json:"user_id"`
}

func (q *Queries) UpdateUserLastAccessedAt(ctx context.Context, arg UpdateUserLastAccessedAtParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserLastAccessedAt, arg.LastAccessedAt, arg.UserID)
}

const updateUserMetamaskAddress = `-- name: UpdateUserMetamaskAddress :execresult
UPDATE users 
SET 
    metamask_address = ?,
    address_changed_at = ?
WHERE user_id = ?
`

type UpdateUserMetamaskAddressParams struct {
	MetamaskAddress  sql.NullString `json:"metamask_address"`
	AddressChangedAt sql.NullTime   `json:"address_changed_at"`
	UserID           string         `json:"user_id"`
}

func (q *Queries) UpdateUserMetamaskAddress(ctx context.Context, arg UpdateUserMetamaskAddressParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserMetamaskAddress, arg.MetamaskAddress, arg.AddressChangedAt, arg.UserID)
}

const updateUserNickname = `-- name: UpdateUserNickname :execresult
UPDATE users 
SET 
    nickname = ?
WHERE user_id = ?
`

type UpdateUserNicknameParams struct {
	Nickname string `json:"nickname"`
	UserID   string `json:"user_id"`
}

func (q *Queries) UpdateUserNickname(ctx context.Context, arg UpdateUserNicknameParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserNickname, arg.Nickname, arg.UserID)
}

const updateUserPhotoURL = `-- name: UpdateUserPhotoURL :execresult
UPDATE users SET photo_url = ?
WHERE user_id = ?
`

type UpdateUserPhotoURLParams struct {
	PhotoUrl sql.NullString `json:"photo_url"`
	UserID   string         `json:"user_id"`
}

func (q *Queries) UpdateUserPhotoURL(ctx context.Context, arg UpdateUserPhotoURLParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserPhotoURL, arg.PhotoUrl, arg.UserID)
}

const updateUserPracBalance = `-- name: UpdateUserPracBalance :execresult
UPDATE users SET prac_balance = ?
WHERE user_id = ?
`

type UpdateUserPracBalanceParams struct {
	PracBalance float64 `json:"prac_balance"`
	UserID      string  `json:"user_id"`
}

func (q *Queries) UpdateUserPracBalance(ctx context.Context, arg UpdateUserPracBalanceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserPracBalance, arg.PracBalance, arg.UserID)
}

const updateUserWmoiBalance = `-- name: UpdateUserWmoiBalance :execresult
UPDATE users SET wmoi_balance = ?
WHERE user_id = ?
`

type UpdateUserWmoiBalanceParams struct {
	WmoiBalance float64 `json:"wmoi_balance"`
	UserID      string  `json:"user_id"`
}

func (q *Queries) UpdateUserWmoiBalance(ctx context.Context, arg UpdateUserWmoiBalanceParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserWmoiBalance, arg.WmoiBalance, arg.UserID)
}

const updateUserWmoiBalanceByRecom = `-- name: UpdateUserWmoiBalanceByRecom :execresult
UPDATE users SET wmoi_balance = ?
WHERE user_id = ?
`

type UpdateUserWmoiBalanceByRecomParams struct {
	WmoiBalance float64 `json:"wmoi_balance"`
	UserID      string  `json:"user_id"`
}

func (q *Queries) UpdateUserWmoiBalanceByRecom(ctx context.Context, arg UpdateUserWmoiBalanceByRecomParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateUserWmoiBalanceByRecom, arg.WmoiBalance, arg.UserID)
}
