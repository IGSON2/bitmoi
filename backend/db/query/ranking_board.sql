-- name: InsertRank :execresult
INSERT INTO ranking_board (
    user_id,
    score_id,
    nickname,
    photo_url,
    comment,
    final_balance
) VALUES (
   ?, ?, ?, ?, ?, ?
);

-- name: GetAllRanks :many
SELECT * FROM ranking_board
ORDER BY balance DESC
LIMIT ?
OFFSET ?;

-- name: GetRankByUserID :one
SELECT * FROM ranking_board
WHERE user_id = ?;

-- name: UpdateUserRank :execresult
UPDATE ranking_board 
SET score_id = ?, final_balance = ?, comment = ?, nickname = ?
WHERE user_id = ?;