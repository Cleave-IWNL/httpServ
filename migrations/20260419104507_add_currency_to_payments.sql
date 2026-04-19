-- +goose Up
ALTER TABLE PAYMENTS ADD COLUMN currency TEXT NOT NULL DEFAULT 'RUB';

-- +goose Down
ALTER TABLE payments DROP COLUMN currency;