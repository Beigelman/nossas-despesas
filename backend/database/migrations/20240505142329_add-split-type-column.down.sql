-- reverse: modify "expenses" table
ALTER TABLE "expenses" DROP COLUMN "split_type";
-- reverse: create enum type "split_type"
DROP TYPE "split_type";
