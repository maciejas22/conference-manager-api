-include .env

MIGRATIONS_PATH = ./migrations

gen:
	protoc --go_out=internal --go_opt=paths=source_relative --go_opt=Mproto/auth.proto=github.com/maciejas22/cm-auth/internal/auth proto/auth.proto
	protoc --go-grpc_out=internal --go-grpc_opt=paths=source_relative --go-grpc_opt=Mproto/auth.proto=github.com/maciejas22/cm-auth/internal/auth proto/auth.proto

start:
	go run ./cmd/main.go
build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main cmd/main.go

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) force $(version)
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down $(count)
