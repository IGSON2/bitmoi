-- name: InsertCompScore :execresult
INSERT INTO comp_score (
    score_id,
    user_id,
    stage,
    pairname,
    entrytime,
    position,
    leverage,
    outtime,
    entryprice,
    quantity,
    endprice,
    pnl,
    roe
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetCompScore :one
SELECT * FROM comp_score
WHERE user_id = ? AND score_id = ? AND pairname = ?;

-- name: GetCompScoresByScoreID :many
SELECT * FROM comp_score
WHERE score_id = ? AND user_id = ?;

-- name: GetCompScoresByStage :one
SELECT * FROM comp_score
WHERE score_id = ? AND user_id = ? AND stage = ?;

-- name: GetCompScoresByUserID :many
SELECT * FROM comp_score
WHERE user_id = ?
ORDER BY score_id DESC 
LIMIT ?
OFFSET ?;

-- name: GetCompScoreToStage :one
SELECT SUM(pnl) FROM comp_score
WHERE score_id = ? AND user_id = ? AND stage <= ?;

-- name: GetCompStageLenByScoreID :one
SELECT COUNT(stage) FROM comp_score
WHERE score_id = ? AND user_id = ?;

-- name: UpdateCompcScore :execresult
UPDATE comp_score SET pairname = ?, entrytime = ?, outtime = ?, entryprice = ?, endprice = ?, pnl = ?, roe = ?
WHERE user_id = ? AND score_id = ? AND pairname = ?;

-- name: GetUnsettledCompScores :many
SELECT * FROM comp_score
WHERE user_id = ? AND pnl <> 0 AND outtime = 0 AND settled_at IS NULL;

-- name: UpdateCompScoreSettledAt :execresult
UPDATE comp_score SET settled_at = ?
WHERE user_id = ? AND score_id = ?;