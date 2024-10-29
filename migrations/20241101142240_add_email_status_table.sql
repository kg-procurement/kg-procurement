-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_status 
(
    id VARCHAR(15) PRIMARY KEY,
    email_to VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    modified_date timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_status;
-- +goose StatementEnd
