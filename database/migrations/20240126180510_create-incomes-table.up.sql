-- create enum type "income_type"
CREATE TYPE "income_type" AS ENUM ('salary', 'benefit', 'vacation', 'thirteenth_salary', 'other');
-- create "incomes" table
CREATE TABLE "incomes" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "amount_cents" bigint NOT NULL,
  "type" "income_type" NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "income_user_id_fk" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
