-- name: GetAdminUsers :many
SELECT id, user_id, nickname, prac_balance, wmoi_balance, recommender_code, created_at, last_accessed_at,
(SELECT COUNT(*) FROM accumulation_history WHERE users.user_id = accumulation_history.to_user AND accumulation_history.title='출석 체크 보상') AS attendance,
(SELECT COUNT(*) FROM recommend_history WHERE users.user_id = recommend_history.recommender) AS referral,
(SELECT COUNT(*) FROM prac_score WHERE users.user_id = prac_score.user_id AND prac_score.pnl >= 0) AS prac_win,
(SELECT COUNT(*) FROM prac_score WHERE users.user_id = prac_score.user_id AND prac_score.pnl < 0) AS prac_lose,
(SELECT COUNT(*) FROM comp_score WHERE users.user_id = comp_score.user_id AND comp_score.pnl >= 0) AS comp_win,
(SELECT COUNT(*) FROM comp_score WHERE users.user_id = comp_score.user_id AND comp_score.pnl < 0) AS comp_lose
FROM users
LIMIT ? OFFSET ?;

-- name: GetAdminPracScores :many
SELECT u.id, u.nickname, u.user_id, p.quantity, p.entryprice, p.position, p.leverage, p.roe, p.pnl, p.entrytime, p.outtime, p.settled_at ,p.created_at
FROM users u LEFT JOIN prac_score p on u.user_id = p.user_id;