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
WHERE user_id = ? AND pnl <> 0 AND outtime IS NOT NULL AND settled_at IS NULL;

-- name: UpdateCompScoreSettledAt :execresult
UPDATE comp_score SET settled_at = ?
WHERE user_id = ? AND score_id = ?;

-- name: GetUserCompScoreSummary :one
SELECT 
  SUM(pnl) AS total_pnl,
  COUNT(CASE WHEN  pnl > 0 THEN 1 END) AS total_win,
  COUNT(CASE WHEN s.pnl < 0 THEN 1 END) AS total_lose,
  SUM(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH THEN s.pnl ELSE 0 END) AS monthly_pnl,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH AND s.pnl > 0 THEN 1 END) AS monthly_win,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL 1 MONTH AND s.pnl < 0 THEN 1 END) AS monthly_lose
FROM comp_score s
JOIN users u ON s.user_id = u.user_id
WHERE u.nickname = ?;