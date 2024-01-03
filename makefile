db:
	docker compose up db -d

migrate-diff:
	atlas migrate diff $(name) -c file://database/atlas.hcl --env local

migrate-create:
	atlas migrate new $(name) -c file://database/atlas.hcl --env local

migrate-apply:
	atlas migrate apply -c file://database/atlas.hcl --env local

migrate-up:
	migrate -path "./database/migrations" -database "postgres://luda:luda@localhost:5432/app?sslmode=disable" up

dev:
	ENV=development go run main.go

unit:
		go test -json -v $$(go list ./... | grep -e handler -e usecase -e pkg | grep -v mocks ) 2>&1 | tee /tmp/gotest.log | gotestfmt

integration:
		export MIGRATIONS_PATH="file://$(shell pwd)/database/migrations"; \
		go test -json -v $$(go list ./... | grep -e postgres) 2>&1 | tee /tmp/gotest.log | gotestfmt

e2e:
		export MIGRATIONS_PATH="file://$(shell pwd)/database/migrations"; \
		go test -json -v $$(go list ./... | grep -e e2e) 2>&1 | tee /tmp/gotest.log | gotestfmt


test: unit integration e2e
