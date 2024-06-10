-- name: CreateRecommendHistory :execresult
INSERT INTO recommend_history (
    recommender,
    new_member
) VALUES (
    ?, ?
);