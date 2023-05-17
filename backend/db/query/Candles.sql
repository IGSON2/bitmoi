-- name: InsertCandles :execresult
INSERT INTO candles_4h (
    name,
    open,
    close,
    high,
    low,
    time,
    volume,
    color
) VALUES (
  ?,?,?,?,?,?,?,?
);

-- name: GetOneCandle :one
SELECT * FROM candles_4h
WHERE name = ? AND time = ?;

-- name: GetCandles :many
SELECT * FROM candles_4h 
WHERE name = ? 
ORDER BY time DESC 
LIMIT ?;