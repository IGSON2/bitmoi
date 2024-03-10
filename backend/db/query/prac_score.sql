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
    roe
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetPracScore :one
SELECT * FROM prac_score
WHERE user_id = ? AND score_id = ? AND pairname = ?;

-- name: GetPracScoresByStage :one
SELECT * FROM prac_score
WHERE score_id = ? AND user_id = ? AND stage = ?;

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
UPDATE prac_score SET outtime = ?, endprice = ?, pnl = ?, roe = ?
WHERE user_id = ? AND score_id = ? AND stage = ?;

-- name: GetUnsettledPracScores :many
SELECT * FROM prac_score
WHERE user_id = ? AND pnl <> 0 AND outtime = 0 AND settled_at IS NULL;

-- name: UpdatePracScoreSettledAt :execresult
UPDATE prac_score SET settled_at = ?
WHERE user_id = ? AND score_id = ?;

-- name: GetUserPracScoreSummary :one
SELECT 
  SUM(pnl) AS total_pnl,
  COUNT(CASE WHEN s.pnl >= 0 THEN 1 END) AS total_win,
  COUNT(CASE WHEN s.pnl < 0 THEN 1 END) AS total_lose,
  SUM(CASE WHEN s.created_at >= CURDATE() - INTERVAL WEEKDAY(CURDATE()) DAY THEN s.pnl ELSE 0 END) AS weekly_pnl,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL WEEKDAY(CURDATE()) DAY AND s.pnl >= 0 THEN 1 END) AS weekly_win,
  COUNT(CASE WHEN s.created_at >= CURDATE() - INTERVAL WEEKDAY(CURDATE()) DAY AND s.pnl < 0 THEN 1 END) AS weekly_lose
FROM prac_score s
JOIN users u ON s.user_id = u.user_id
WHERE u.nickname = ?;

-- name: GetUserPracRankByPNL :many
SELECT users.nickname, ranks.sum AS "sum" FROM 
(SELECT user_id,SUM(pnl) AS "sum" FROM prac_score WHERE prac_score.created_at >= ? AND prac_score.created_at < ? GROUP BY user_id) ranks
JOIN users ON ranks.user_id = users.user_id
ORDER BY ranks.sum DESC LIMIT ? OFFSET ?;

-- name: GetUserPracRankByROE :many
SELECT users.nickname, ranks.sum AS "sum" FROM 
(SELECT user_id,SUM(roe) AS "sum" FROM prac_score WHERE prac_score.created_at >= ? AND prac_score.created_at < ? GROUP BY user_id) ranks
JOIN users ON ranks.user_id = users.user_id
ORDER BY ranks.sum DESC LIMIT ? OFFSET ?;