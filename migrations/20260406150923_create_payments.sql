-- +goose Up
CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY,
    amount INTEGER NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS payments;
