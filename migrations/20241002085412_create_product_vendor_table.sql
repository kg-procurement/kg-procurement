-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_vendor
(
    product_id VARCHAR(15),
    vendor_id  VARCHAR(15),
    PRIMARY KEY (product_id, vendor_id),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (vendor_id) REFERENCES vendor (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_vendor;
-- +goose StatementEnd
