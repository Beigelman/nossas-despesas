-- reverse: create "expenses_latest" view
DROP INDEX "idx_expenses_latest_group_created";
DROP INDEX "idx_expenses_latest_created_at";
DROP INDEX "idx_expenses_latest_group_deleted";
DROP INDEX "idx_expenses_id_version";
DROP VIEW "expenses_latest";
