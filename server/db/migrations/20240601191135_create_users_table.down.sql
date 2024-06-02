-- +goose Up
-- +goose StatementBegin

-- This section intentionally left blank.

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd