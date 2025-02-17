module github.com/maciejas22/conference-manager-api/cm-conferences

go 1.23.1

require (
	github.com/jackc/pgx/v5 v5.7.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/maciejas22/conference-manager-api/cm-proto v0.0.0
	github.com/stripe/stripe-go/v81 v81.0.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)

replace github.com/maciejas22/conference-manager-api/cm-proto => ../cm-proto

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
)
