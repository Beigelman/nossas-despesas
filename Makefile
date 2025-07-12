# docker
db:
	docker compose up db -d
# Run
dev: db
	ENV=development go run main.go



# Migrations
migrate-diff:
	atlas migrate diff $(name) -c file://database/atlas.hcl --env local

migrate-new:
	atlas migrate new $(name) -c file://database/atlas.hcl --env local

migrate-hash:
	 atlas migrate hash -c file://database/atlas.hcl --env local

migrate-up:
	./database/migrate.sh up "./database/migrations"

migrate-down:
	./database/migrate.sh down "./database/migrations"

migrate-force:
	migrate -path "./database/migrations" -database "postgres://root:root@localhost:5432/app?sslmode=disable" force $(version)

# Tests and format
format:
	goimports -w -l ./internal
	golangci-lint run --fix

mock:
	mockery

unit:
		go test -v $$(go list ./internal/modules/... | grep -e handler -e usecase -e pkg)

integration:
		go test -v $$(go list ./internal/modules/... | grep -e postgres)

test:
		go test -json -cover ./internal/modules/... | tparse

# Scripts
create-users:
	go run ./scripts/main.go create-users

import-split:
	go run ./scripts/main.go import-from-split-wize

import-incomes:
	go run ./scripts/main.go import-incomes

reset-app: db migrate-down migrate-up create-users import-incomes import-split


