-- name: CreateWmoiTransaction :execresult
INSERT INTO wmoi_transaction (
    user_id,
    to,
    amount,
    title
) VALUES (
    ?, ?, ?, ?
);