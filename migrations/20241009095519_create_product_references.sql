-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_category
(
    id              VARCHAR(15) PRIMARY KEY,
    name            VARCHAR(127),
    code            VARCHAR(15),
    description     VARCHAR(255),
    parent_id       VARCHAR(15),
    specialist_bpid VARCHAR(15),
    modified_date   timestamp,
    modified_by     VARCHAR(127)
);

CREATE TABLE product_type
(
    id            VARCHAR(15) PRIMARY KEY,
    name          VARCHAR(127),
    description   VARCHAR(255),
    goods         VARCHAR(15),
    asset         VARCHAR(15),
    stock         VARCHAR(15),
    modified_date timestamp,
    modified_by   VARCHAR(127)
);

CREATE TABLE uom
(
    id            VARCHAR(15) PRIMARY KEY,
    name          VARCHAR(127),
    description   VARCHAR(255),
    dimension     VARCHAR(15),
    sap_code      VARCHAR(15),
    modified_date timestamp,
    modified_by   VARCHAR(127)
);

ALTER TABLE product
    ADD CONSTRAINT fk_product_category
        FOREIGN KEY (product_category_id) REFERENCES product_category (id)
            ON DELETE SET NULL,
    ADD CONSTRAINT fk_uom
        FOREIGN KEY (uom_id) REFERENCES uom (id)
        ON DELETE SET NULL,
    ADD CONSTRAINT fk_product_type
        FOREIGN KEY (product_type_id) REFERENCES product_type (id)
        ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product
DROP CONSTRAINT IF EXISTS fk_product_category,
    DROP CONSTRAINT IF EXISTS fk_uom,
    DROP CONSTRAINT IF EXISTS fk_product_type;

DROP TABLE product_category;
DROP TABLE product_type;
DROP TABLE uom;
-- +goose StatementEnd
