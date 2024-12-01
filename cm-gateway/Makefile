-include .env
-include .env.prod
-include .env.dev

MIGRATIONS_PATH = ./migrations

gen:
	go run github.com/99designs/gqlgen generate

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main cmd/cm-gateway/main.go

start:
	go run cmd/conference-manager-api/main.go

start-prod:
	./bin/conference-manager-api

lint:
	golangci-lint run -v

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down $(count)

clean:
	rm -rf bin/

.PHONY: gen start lint migrate-create migrate-up migrate-down build start-prod clean
