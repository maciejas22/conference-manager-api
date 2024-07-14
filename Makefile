include .env

MIGRATIONS_PATH = db/migrations

gen:
	go run github.com/99designs/gqlgen generate

start:
	go run cmd/conference-manager-api/main.go 

lint:
	golangci-lint run -v

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down $(count)
