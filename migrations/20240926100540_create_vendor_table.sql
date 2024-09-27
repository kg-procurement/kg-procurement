-- +goose Up
-- +goose StatementBegin
CREATE TABLE "vendor" (
  "id" int PRIMARY KEY,
  "name" varchar(255),
  "bp_id" int,
  "bp_name" varchar(255),
  "rating"  int,
  "area_group_id" int,
  "area_group_name" varchar(255),
  "sap_code"	  varchar(255),
  "modified_date" timestamp,
  "modified_by" int,
  "dt" date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "vendor";
-- +goose StatementEnd
