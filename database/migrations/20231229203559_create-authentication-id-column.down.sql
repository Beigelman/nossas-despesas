-- reverse: create index "authentication_id_idx" to table: "users"
DROP INDEX "authentication_id_idx";
-- reverse: modify "users" table
ALTER TABLE "users" DROP COLUMN "authentication_id";
