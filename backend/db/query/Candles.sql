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

-- name: Get1dVolSumPriceAVG :one
SELECT SUM(volume) AS volsum, AVG(close) AS priceavg FROM candles_1d WHERE name = ? AND time <= ?;

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

-- name: Get4hVolSumPriceAVG :one
SELECT SUM(volume) AS volsum, AVG(close) AS priceavg FROM candles_4h WHERE name = ? AND time <= ?;

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

-- name: Get1hEntryTimestamp :one
SELECT time FROM candles_1h 
WHERE name = ?  AND time <= ?
ORDER BY time desc 
LIMIT 1;

-- name: Get1hResult :many
SELECT * FROM candles_1h 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get1hMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_1h
WHERE name = ?;

-- name: Get1hVolSumPriceAVG :one
SELECT SUM(volume) AS volsum, AVG(close) AS priceavg FROM candles_1h WHERE name = ? AND time <= ?;

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

-- name: Get15mVolSumPriceAVG :one
SELECT SUM(volume) AS volsum, AVG(close) AS priceavg FROM candles_15m WHERE name = ? AND time <= ?;

-- name: Insert5mCandles :execresult
INSERT INTO candles_5m (
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

-- name: GetOne5mCandle :one
SELECT * FROM candles_5m
WHERE name = ? AND time = ?;

-- name: Get5mCandles :many
SELECT * FROM candles_5m 
WHERE name = ?  AND time <= ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get5mResult :many
SELECT * FROM candles_5m 
WHERE name = ? AND time > ?
ORDER BY time ASC 
LIMIT ?;

-- name: Get5mMinMaxTime :one
SELECT MIN(time), MAX(time)
FROM candles_5m
WHERE name = ?;

-- name: Get5mVolSumPriceAVG :one
SELECT SUM(volume) AS volsum, AVG(close) AS priceavg FROM candles_5m WHERE name = ? AND time <= ?;

-- name: GetAllParisInDB :many
SELECT DISTINCT name from candles_1h
ORDER BY name;