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

-- name: GetAdminScores :many
SELECT u.id, u.nickname, p.* ,a.min_roe, a.max_roe, a.after_outtime from prac_score p
 INNER JOIN users u ON p.user_id = u.user_id
 INNER JOIN prac_after_score a ON p.user_id = a.user_id AND p.score_id = a.score_id
 LIMIT ? OFFSET ?;