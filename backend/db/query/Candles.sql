-- name: Insert1dCandles :execresult
INSERT INTO candles_1d (
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

-- name: GetOne1dCandle :one
SELECT * FROM candles_1d
WHERE name = ? AND time = ?;

-- name: Get1dCandles :many
SELECT * FROM candles_1d 
WHERE name = ? AND time <= ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get1dResult :many
SELECT * FROM candles_1d 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get1dMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_1d
WHERE name = ?;

-- name: Insert4hCandles :execresult
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

-- name: GetOne4hCandle :one
SELECT * FROM candles_4h
WHERE name = ? AND time = ?;

-- name: Get4hCandles :many
SELECT * FROM candles_4h 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get4hResult :many
SELECT * FROM candles_4h 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get4hMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_4h
WHERE name = ?;


-- name: Insert1hCandles :execresult
INSERT INTO candles_1h (
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

-- name: GetOne1hCandle :one
SELECT * FROM candles_1h
WHERE name = ? AND time = ?;

-- name: Get1hCandles :many
SELECT * FROM candles_1h 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get1hResult :many
SELECT * FROM candles_1h 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get1hMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_1h
WHERE name = ?;

-- name: Insert15mCandles :execresult
INSERT INTO candles_15m (
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

-- name: GetOne15mCandle :one
SELECT * FROM candles_15m
WHERE name = ? AND time = ?;

-- name: Get15mCandles :many
SELECT * FROM candles_15m 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get15mResult :many
SELECT * FROM candles_15m 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get15mMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_15m
WHERE name = ?;