-- create "expenses_latest" view
CREATE VIEW "expenses_latest" (
  "id",
  "name",
  "amount_cents",
  "refund_amount_cents",
  "description",
  "group_id",
  "category_id",
  "split_ratio",
  "split_type",
  "payer_id",
  "receiver_id",
  "document_search",
  "created_at",
  "updated_at",
  "deleted_at",
  "version"
) AS SELECT DISTINCT ON (expenses.id) expenses.id,
    expenses.name,
    expenses.amount_cents,
    expenses.refund_amount_cents,
    expenses.description,
    expenses.group_id,
    expenses.category_id,
    expenses.split_ratio,
    expenses.split_type,
    expenses.payer_id,
    expenses.receiver_id,
    expenses.document_search,
    expenses.created_at,
    expenses.updated_at,
    expenses.deleted_at,
    expenses.version
   FROM expenses
  ORDER BY expenses.id DESC, expenses.version DESC;

-- create indexes to optimize the view performance
CREATE INDEX "idx_expenses_id_version" ON "expenses" ("id" DESC, "version" DESC);
CREATE INDEX "idx_expenses_latest_group_deleted" ON "expenses" ("group_id", "deleted_at") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_expenses_latest_created_at" ON "expenses" ("created_at");
CREATE INDEX "idx_expenses_latest_group_created" ON "expenses" ("group_id", "created_at");
