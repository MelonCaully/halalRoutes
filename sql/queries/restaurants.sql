-- name: CreateRestaurant :exec
INSERT INTO halal_restaurants (id, name, region, address, website, source, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
);