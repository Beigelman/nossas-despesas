-- create enum type "authentication_type"
CREATE TYPE "authentication_type" AS ENUM ('credentials', 'google');
-- create enum type "income_type"
CREATE TYPE "income_type" AS ENUM ('salary', 'benefit', 'vacation', 'thirteenth_salary', 'other');
-- create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "profile_picture" character varying(255) NULL,
  "group_id" bigint NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "email_unique_idx" to table: "users"
CREATE UNIQUE INDEX "email_unique_idx" ON "users" ("email");
-- create "authentications" table
CREATE TABLE "authentications" (
  "id" bigserial NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(255) NULL,
  "provider_id" character varying(255) NULL,
  "type" "authentication_type" NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "auth_email_fk" FOREIGN KEY ("email") REFERENCES "users" ("email") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "email_type_unique_idx" to table: "authentications"
CREATE UNIQUE INDEX "email_type_unique_idx" ON "authentications" ("email", "type");
-- create "category_groups" table
CREATE TABLE "category_groups" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "icon" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create "categories" table
CREATE TABLE "categories" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "icon" character varying(255) NOT NULL,
  "category_group_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "category_group_id_fk" FOREIGN KEY ("category_group_id") REFERENCES "category_groups" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "groups" table
CREATE TABLE "groups" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create "expenses" table
CREATE TABLE "expenses" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "amount_cents" bigint NOT NULL,
  "refund_amount_cents" bigint NULL,
  "description" character varying(255) NOT NULL,
  "group_id" bigint NOT NULL,
  "category_id" bigint NOT NULL,
  "split_ratio" jsonb NOT NULL,
  "payer_id" bigint NOT NULL,
  "receiver_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id", "version"),
  CONSTRAINT "category_id_fk" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "group_id_fk" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "payer_id_fk" FOREIGN KEY ("payer_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "receiver_id_fk" FOREIGN KEY ("receiver_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
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
