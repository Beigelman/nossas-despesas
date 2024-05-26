-- modify "expenses" table
ALTER TABLE "expenses" ADD COLUMN "document_search" tsvector NOT NULL GENERATED ALWAYS AS (setweight(to_tsvector('portuguese'::regconfig, (COALESCE(name, ''::character varying))::text), 'A'::"char") || setweight(to_tsvector('portuguese'::regconfig, (COALESCE(description, ''::character varying))::text), 'B'::"char")) STORED;
-- create index "document_search_idx" to table: "expenses"
CREATE INDEX "document_search_idx" ON "expenses" USING gin ("document_search");
