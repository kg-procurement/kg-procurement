-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS product_vendor;
CREATE TABLE product_vendor
(
    id                    VARCHAR(15) PRIMARY KEY,
    product_id            VARCHAR(15),
    code                  VARCHAR(15),
    name                  VARCHAR(255),
    income_tax_id         VARCHAR(15),
    income_tax_name       VARCHAR(255),
    income_tax_percentage VARCHAR(15),
    description           VARCHAR(255),
    uom_id                VARCHAR(15),
    sap_code              VARCHAR(15),
    modified_date         TIMESTAMP,
    modified_by           VARCHAR(15),
    FOREIGN KEY (product_id) REFERENCES product (id),
    FOREIGN KEY (uom_id) REFERENCES uom (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_vendor;
-- +goose StatementEnd
