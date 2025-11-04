-- reverse: create index "document_search_idx" to table: "expenses"
DROP INDEX "document_search_idx";
-- reverse: modify "expenses" table
ALTER TABLE "expenses" DROP COLUMN "document_search";
