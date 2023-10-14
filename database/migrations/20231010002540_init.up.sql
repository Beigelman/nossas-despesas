-- create "categories" table
CREATE TABLE "categories" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "icon" character varying(255) NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create "groups" table
CREATE TABLE "groups" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create "users" table
CREATE TABLE "users" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "profile_picture" character varying(255) NULL,
  "group_id" integer NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create "expenses" table
CREATE TABLE "expenses" (
  "id" serial NOT NULL,
  "name" character varying(255) NOT NULL,
  "amount_cents" bigint NOT NULL,
  "description" character varying(255) NOT NULL,
  "group_id" integer NOT NULL,
  "category_id" integer NOT NULL,
  "split_ratio" jsonb NOT NULL,
  "payer_id" integer NOT NULL,
  "receiver_id" integer NOT NULL,
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
