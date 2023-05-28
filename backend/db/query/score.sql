-- name: InsertScore :execresult
INSERT INTO score (
    score_id,
    user_id,
    stage,
    pairname,
    entrytime,
    position,
    leverage,
    outtime,
    entryprice,
    endprice,
    pnl,
    roe
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetScore :one
SELECT * FROM score
WHERE score_id = ? AND stage = ?;

-- name: GetScoresByScoreID :many
SELECT * FROM score
WHERE score_id = ? AND user_id = ?;

-- name: GetScoresByUserID :many
SELECT * FROM score
WHERE user_id = ?
ORDER BY score_id DESC 
LIMIT ?
OFFSET ?;

-- name: GetScoreToStage :one
SELECT SUM(pnl) FROM score
WHERE score_id = ? AND user_id = ? AND stage <= ?;

-- name: GetStageLenByScoreID :one
SELECT COUNT(stage) FROM score
WHERE score_id = ? AND user_id = ?;