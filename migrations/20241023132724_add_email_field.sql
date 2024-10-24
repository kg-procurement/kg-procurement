-- +goose Up
-- +goose StatementBegin
ALTER TABLE vendor
    ADD email VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vendor
    DROP COLUMN email;
-- +goose StatementEnd
