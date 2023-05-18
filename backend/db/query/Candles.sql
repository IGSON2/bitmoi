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
WHERE name = ? AND id = ?;

-- name: Get1dCandles :many
SELECT * FROM candles_1d 
WHERE name = ? 
ORDER BY time ASC 
LIMIT ?;

-- name: Counting1dCandles :one
SELECT count(time) FROM candles_1d;

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
WHERE name = ? AND id = ?;

-- name: Get4hCandles :many
SELECT * FROM candles_4h 
WHERE name = ? 
ORDER BY time ASC 
LIMIT ?;

-- name: Counting4hCandles :one
SELECT count(time) FROM candles_4h;

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
WHERE name = ? AND id = ?;

-- name: Get1hCandles :many
SELECT * FROM candles_1h 
WHERE name = ? 
ORDER BY time ASC 
LIMIT ?;

-- name: Counting1hCandles :one
SELECT count(time) FROM candles_1h;

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
WHERE name = ? AND id = ?;

-- name: Get15mCandles :many
SELECT * FROM candles_15m 
WHERE name = ? 
ORDER BY time ASC 
LIMIT ?;

-- name: Counting15mCandles :one
SELECT count(time) FROM candles_15m;