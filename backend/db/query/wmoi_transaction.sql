-- name: CreateWmoiTransaction :execresult
INSERT INTO wmoi_transaction (
    from_user,
    to_user,
    amount,
    title
) VALUES (
    ?, ?, ?, ?
);