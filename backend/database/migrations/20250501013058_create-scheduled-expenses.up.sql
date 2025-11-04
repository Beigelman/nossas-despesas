-- create "scheduled_expenses" table
CREATE TABLE "scheduled_expenses" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "amount_cents" bigint NOT NULL,
  "description" text NOT NULL,
  "group_id" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "split_type" "split_type" NOT NULL,
  "payer_id" bigint NOT NULL,
  "receiver_id" bigint NOT NULL,
  "frequency_in_days" integer NOT NULL,
  "last_generated_at" date NULL,
  "is_active" boolean NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
