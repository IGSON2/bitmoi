-- name: InsertPracAfterScore :execresult
INSERT INTO prac_after_score (
    score_id,
    user_id,
    max_roe,
    min_roe,
    after_outtime
) VALUES (
    ?, ?, ?, ?, ?
);