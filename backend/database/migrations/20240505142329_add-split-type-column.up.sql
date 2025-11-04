-- create enum type "split_type"
CREATE TYPE "split_type" AS ENUM ('equal', 'proportional', 'transfer');
-- modify "expenses" table
ALTER TABLE "expenses" ADD COLUMN "split_type" "split_type" NOT NULL;
