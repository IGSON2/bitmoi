-- name: CreateAccumulationHist :execresult
INSERT INTO accumulation_history (
    to_user,
    amount,
    title
) VALUES (
    ?, ?, ?
);

-- name: GetAccumulationHist :many
SELECT * FROM accumulation_history
WHERE to_user = ? AND title LIKE ?
ORDER BY created_at DESC 
LIMIT ?
OFFSET ?;