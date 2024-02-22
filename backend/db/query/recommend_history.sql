-- name: CreateRecommendHistory :execresult
INSERT INTO recommend_history (
    from_user,
    to_user
) VALUES (
    ?, ?
);