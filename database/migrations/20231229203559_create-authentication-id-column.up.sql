-- modify "users" table
ALTER TABLE "users" ADD COLUMN "authentication_id" character varying(255) NULL;
-- create index "authentication_id_idx" to table: "users"
CREATE UNIQUE INDEX "authentication_id_idx" ON "users" ("authentication_id");
