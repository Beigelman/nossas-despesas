-- create enum type "group_invites_status"
CREATE TYPE "group_invites_status" AS ENUM ('pending', 'sent', 'accepted');
-- create "group_invites" table
CREATE TABLE "group_invites" (
  "id" bigserial NOT NULL,
  "group_id" bigint NOT NULL,
  "email" character varying(255) NOT NULL,
  "token" uuid NOT NULL DEFAULT gen_random_uuid(),
  "status" "group_invites_status" NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "deleted_at" timestamptz NULL,
  "version" integer NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "group_invite_email_idx" to table: "group_invites"
CREATE INDEX "group_invite_email_idx" ON "group_invites" ("email");
-- create index "group_invite_token_idx" to table: "group_invites"
CREATE UNIQUE INDEX "group_invite_token_idx" ON "group_invites" ("token");
