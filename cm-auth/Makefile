-include .env

MIGRATIONS_PATH = ./migrations

start:
	go run ./cmd/main.go
build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main cmd/main.go

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down $(count)
