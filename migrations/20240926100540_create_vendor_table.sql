-- +goose Up
-- +goose StatementBegin
CREATE TABLE "vendor" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "description" varchar,
  "bp_id" int,
  "bp_name" varchar,
  "rating"  int,
  "area_group_id" int,
  "area_group_name" varchar,
  "sap_code" varchar,
  "modified_date" timestamp,
  "modified_by" int,
  "dt" date
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "vendor";
-- +goose StatementEnd
