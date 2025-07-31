-- +goose Up 
CREATE TABLE crime_stats (
    id UUID PRIMARY KEY,
    neighborhood TEXT NOT NULL,
    total_crime INT NOT NULL,
    violent_crime INT NOT NULL,
    property_crime INT NOT NULL,
    source TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS crime_stats;