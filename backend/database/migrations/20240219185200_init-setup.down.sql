-- reverse: create "incomes" table
DROP TABLE "incomes";
-- reverse: create "expenses" table
DROP TABLE "expenses";
-- reverse: create "groups" table
DROP TABLE "groups";
-- reverse: create "categories" table
DROP TABLE "categories";
-- reverse: create "category_groups" table
DROP TABLE "category_groups";
-- reverse: create index "email_type_unique_idx" to table: "authentications"
DROP INDEX "email_type_unique_idx";
-- reverse: create "authentications" table
DROP TABLE "authentications";
-- reverse: create index "email_unique_idx" to table: "users"
DROP INDEX "email_unique_idx";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create enum type "income_type"
DROP TYPE "income_type";
-- reverse: create enum type "authentication_type"
DROP TYPE "authentication_type";
