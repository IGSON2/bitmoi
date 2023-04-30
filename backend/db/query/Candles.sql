-- name: InsertCandles :execresult
INSERT INTO Candles (
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