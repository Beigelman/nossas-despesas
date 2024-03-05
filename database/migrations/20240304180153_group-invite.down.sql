-- reverse: create index "group_invite_token_idx" to table: "group_invites"
DROP INDEX "group_invite_token_idx";
-- reverse: create index "group_invite_email_idx" to table: "group_invites"
DROP INDEX "group_invite_email_idx";
-- reverse: create "group_invites" table
DROP TABLE "group_invites";
-- reverse: create enum type "group_invites_status"
DROP TYPE "group_invites_status";
