-- name: CreateBiddingHistory :execresult
INSERT INTO bidding_history (
    user_id,
    amount,
    location,
    tx_hash,
    expires_at
) VALUES (
    ?, ?, ?, ?,?
);

-- name: GetHistoryByLocation :many
SELECT * FROM bidding_history 
WHERE location = ? AND expires_at >= now()
ORDER BY amount DESC 
LIMIT ?;

-- name: GetHistoryByUser :many
SELECT * FROM bidding_history 
WHERE user_id = ? AND expires_at >= now()
ORDER BY amount DESC;