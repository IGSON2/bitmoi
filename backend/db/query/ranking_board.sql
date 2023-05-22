-- name: InsertRank :execresult
INSERT INTO ranking_board (
    user_id,
    score_id,
    display_name,
    photo_url,
    comment,
    final_balance
) VALUES (
   ?, ?, ?, ?, ?, ?
);

-- name: GetRankByUserID :one
SELECT * FROM ranking_board
WHERE user_id = ?;

-- name: UpdateUserRank :execresult
UPDATE ranking_board 
SET score_id = ?, final_balance = ?, comment = ?
WHERE user_id = ?;