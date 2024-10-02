-- +goose Up
-- +goose StatementBegin
CREATE TABLE product
(
    id                  VARCHAR(15) PRIMARY KEY,
    product_category_id VARCHAR(15),
    uom_id              VARCHAR(15),
    income_tax_id       VARCHAR(15),
    product_type_id     VARCHAR(15),
    name                VARCHAR(127) NOT NULL,
    description         VARCHAR(255),
    modified_date       timestamp,
    modified_by         VARCHAR(127)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product;
-- +goose StatementEnd
