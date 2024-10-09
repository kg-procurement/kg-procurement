-- +goose Up
-- +goose StatementBegin
CREATE TABLE account (
    id VARCHAR(15) PRIMARY KEY,
    email VARCHAR(63) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account;
-- +goose StatementEnd
