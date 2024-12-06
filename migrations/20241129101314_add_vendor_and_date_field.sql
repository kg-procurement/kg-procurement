-- +goose Up
-- +goose StatementBegin
ALTER TABLE email_status
    ADD date_sent timestamp,
    ADD vendor_id VARCHAR(15);

ALTER TABLE email_status
    ADD CONSTRAINT fk_vendor
        FOREIGN KEY (vendor_id) REFERENCES vendor (id)
        ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE email_status
    DROP CONSTRAINT IF EXISTS fk_vendor;

ALTER TABLE email_status
    DROP COLUMN IF EXISTS vendor_id,
    DROP COLUMN IF EXISTS date_sent;
-- +goose StatementEnd
