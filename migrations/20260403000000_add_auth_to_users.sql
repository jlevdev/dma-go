-- +goose Up
ALTER TABLE Users
    ADD COLUMN GoogleId TEXT UNIQUE,
    ADD COLUMN Email TEXT UNIQUE,
    ADD COLUMN DisplayName TEXT;

-- +goose Down
ALTER TABLE Users
    DROP COLUMN GoogleId,
    DROP COLUMN Email,
    DROP COLUMN DisplayName;
