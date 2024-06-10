-- name: CreateUsedToken :execresult
INSERT INTO used_token (
    score_id,
    user_id,
    metamask_address
) VALUES (
    ?, ?, ?
);