-- name: CreateWmoiMintinghist :execresult
INSERT INTO wmoi_minting_history (
    to_user,
    amount,
    title
) VALUES (
    ?, ?, ?
);

-- name: GetWmoiMintingHist :many
SELECT * FROM wmoi_minting_history
WHERE to_user = ? AND title LIKE ?
ORDER BY created_at DESC 
LIMIT ?
OFFSET ?;