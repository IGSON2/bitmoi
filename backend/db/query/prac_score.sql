-- name: InsertPracScore :execresult
INSERT INTO prac_score (
    score_id,
    user_id,
    stage,
    pairname,
    entrytime,
    position,
    leverage,
    outtime,
    quantity,
    entryprice,
    endprice,
    pnl,
    roe,
    remain_balance
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetPracScore :one
SELECT * FROM prac_score
WHERE user_id = ? AND score_id = ? AND stage = ?;

-- name: GetPracScoresByScoreID :many
SELECT * FROM prac_score
WHERE score_id = ? AND user_id = ?;

-- name: GetPracScoresByUserID :many
SELECT * FROM prac_score
WHERE user_id = ?
ORDER BY score_id DESC 
LIMIT ?
OFFSET ?;

-- name: GetPracScoreToStage :one
SELECT SUM(pnl) FROM prac_score
WHERE score_id = ? AND user_id = ? AND stage <= ?;

-- name: GetPracStageLenByScoreID :one
SELECT COUNT(stage) FROM prac_score
WHERE score_id = ? AND user_id = ?;

-- name: UpdatePracScore :execresult
UPDATE prac_score SET outtime = ?, endprice = ?, pnl = ?, roe = ?, remain_balance = ?
WHERE user_id = ? AND score_id = ? AND stage = ?;