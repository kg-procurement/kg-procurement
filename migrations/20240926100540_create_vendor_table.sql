-- +goose Up
-- +goose StatementBegin
CREATE TABLE "vendor"
(
    "id"              varchar(15) PRIMARY KEY,
    "name"            varchar(127),
    "description"     varchar(127),
    "bp_id"           varchar(15),
    "bp_name"         varchar(127),
    "rating"          int,
    "area_group_id"   varchar(15),
    "area_group_name" varchar(127),
    "sap_code"        varchar(127),
    "modified_date"   timestamp,
    "modified_by"     VARCHAR(127),
    "dt"              date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "vendor";
-- +goose StatementEnd
