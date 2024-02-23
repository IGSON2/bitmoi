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

-- name: Get1dCandles :many
SELECT * FROM candles_1d 
WHERE name = ? AND time <= ?
ORDER BY time DESC 
LIMIT ?;

-- name: Get1dCandlesRnage :many
SELECT * FROM candles_1d 
WHERE name = ? AND time > ? AND time <= ?
ORDER BY time DESC;

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

-- name: DeletePairs1d :execresult
DELETE FROM candles_1d WHERE name = ?;

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

-- name: Get4hCandles :many
SELECT * FROM candles_4h 
WHERE name = ?  AND time <= ?
ORDER BY time DESC 
LIMIT ?;

-- name: Get4hCandlesRnage :many
SELECT * FROM candles_4h 
WHERE name = ? AND time > ? AND time <= ?
ORDER BY time DESC;

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

-- name: DeletePairs4h :execresult
DELETE FROM candles_4h WHERE name = ?;

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

-- name: Get1hCandles :many
SELECT * FROM candles_1h 
WHERE name = ?  AND time <= ?
ORDER BY time DESC 
LIMIT ?;

-- name: Get1hCandlesRnage :many
SELECT * FROM candles_1h 
WHERE name = ? AND time > ? AND time <= ?
ORDER BY time DESC;

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

-- name: DeletePairs1h :execresult
DELETE FROM candles_1h WHERE name = ?;

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

-- name: Get15mCandles :many
SELECT * FROM candles_15m 
WHERE name = ?  AND time <= ?
ORDER BY time DESC 
LIMIT ?;

-- name: Get15mCandlesRnage :many
SELECT * FROM candles_15m 
WHERE name = ? AND time > ? AND time <= ?
ORDER BY time DESC;

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

-- name: DeletePairs15m :execresult
DELETE FROM candles_15m WHERE name = ?;

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

-- name: Get5mCandles :many
SELECT * FROM candles_5m 
WHERE name = ?  AND time <= ?
ORDER BY time DESC 
LIMIT ?;

-- name: Get5mCandlesRnage :many
SELECT * FROM candles_5m 
WHERE name = ? AND time > ? AND time <= ?
ORDER BY time DESC;

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

-- name: DeletePairs5m :execresult
DELETE FROM candles_5m WHERE name = ?;

-- --------utils----------------

-- name: GetAllPairsInDB1H :many
SELECT DISTINCT name from candles_1h
ORDER BY name;

-- name: GetAllPairsInDB1D :many
SELECT DISTINCT name from candles_1d
ORDER BY name;

-- name: GetUnder1YPairs :many
SELECT name
FROM candles_1d
GROUP BY name
HAVING COUNT(name) < 365;