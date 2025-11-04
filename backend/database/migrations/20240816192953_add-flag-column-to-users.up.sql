-- modify "users" table
ALTER TABLE "users" ADD COLUMN "flags" text[] NOT NULL DEFAULT ARRAY[]::text[];
