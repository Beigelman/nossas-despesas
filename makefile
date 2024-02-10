db:
	docker compose up db -d

migrate-diff:
	atlas migrate diff $(name) -c file://database/atlas.hcl --env local

migrate-new:
	atlas migrate new $(name) -c file://database/atlas.hcl --env local

migrate-hash:
	 atlas migrate hash -c file://database/atlas.hcl --env local

migrate-up:
	migrate -path "./database/migrations" -database "postgres://root:root@localhost:5432/app?sslmode=disable" up

migrate-down:
	migrate -path "./database/migrations" -database "postgres://root:root@localhost:5432/app?sslmode=disable" down

migrate-force:
	migrate -path "./database/migrations" -database "postgres://root:root@localhost:5432/app?sslmode=disable" force $(version)

dev:
	ENV=development go run main.go

mock:
	mockery

unit:
		go test -json -v $$(go list ./... | grep -e handler -e usecase -e pkg | grep -v mocks ) 2>&1 | tee /tmp/gotest.log | gotestfmt

integration:
		export MIGRATIONS_PATH="file://$(shell pwd)/database/migrations"; \
		go test -json -v $$(go list ./... | grep -e postgres) 2>&1 | tee /tmp/gotest.log | gotestfmt

e2e:
		export MIGRATIONS_PATH="file://$(shell pwd)/database/migrations"; \
		go test -json -v $$(go list ./... | grep -e e2e) 2>&1 | tee /tmp/gotest.log | gotestfmt


test: unit integration e2e

import-split:
	go run ./scripts/main.go import-from-split-wize -d 1 -l 2 -g 1

create-users:
	go run ./scripts/main.go create-users

reset-app: db migrate-down migrate-up create-users import-split


