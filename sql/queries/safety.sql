-- name: CreateCrimeStats :one
INSERT INTO crime_stats (id, neighborhood, total_crime, violent_crime, property_crime, source, created_at, updated_at)
VALUES(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
)
RETURNING *;

-- name: GetAllCrimeStats :many
SELECT * FROM crime_stats ORDER BY neighborhood;