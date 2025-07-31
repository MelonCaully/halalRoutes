-- +goose Up 
CREATE TABLE neighborhood_safety (
    id UUID PRIMARY KEY,
    neighborhood TEXT NOT NULL,
    safety_score INT,
    source TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS neighborhood_safety;