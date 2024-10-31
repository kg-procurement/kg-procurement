-- +goose Up
-- +goose StatementBegin
CREATE TABLE price
(
    id                   VARCHAR(15) PRIMARY KEY,
    purchasing_org_id    VARCHAR(15),
    purchasing_org_name  VARCHAR(255),
    vendor_id            VARCHAR(15),
    product_vendor_id    VARCHAR(15),
    quantity_min         INT,
    quantity_max         INT,
    quantity_uom_id      VARCHAR(15),
    lead_time_min        INT,
    lead_time_max        INT,
    currency_id          VARCHAR(15),
    currency_name        VARCHAR(255),
    currency_code        VARCHAR(15),
    price                NUMERIC(15, 2),
    price_quantity       INT,
    price_uom_id         VARCHAR(15),
    valid_from           TIMESTAMP,
    valid_to             TIMESTAMP,
    valid_pattern_id     VARCHAR(15),
    valid_pattern_name   VARCHAR(255),
    area_group_id        VARCHAR(15),
    area_group_name      VARCHAR(255),
    reference_number     VARCHAR(127),
    reference_date       DATE,
    document_type_id     VARCHAR(15),
    document_type_name   VARCHAR(255),
    document_id          VARCHAR(15),
    item_id              VARCHAR(15),
    term_of_payment_id   VARCHAR(15),
    term_of_payment_days INT,
    term_of_payment_text VARCHAR(255),
    invocation_order     INT,
    modified_date        TIMESTAMP,
    modified_by          VARCHAR(255),

    CONSTRAINT fk_vendor FOREIGN KEY (vendor_id) REFERENCES vendor (id),
    CONSTRAINT fk_product_vendor FOREIGN KEY (product_vendor_id) REFERENCES product_vendor (id),
    CONSTRAINT fk_quantity_uom FOREIGN KEY (quantity_uom_id) REFERENCES uom (id),
    CONSTRAINT fk_price_uom FOREIGN KEY (price_uom_id) REFERENCES uom (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE price;
-- +goose StatementEnd
