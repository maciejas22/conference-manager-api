# builder 
FROM golang:latest AS builder

WORKDIR /app/cm-auth

COPY cm-auth/go.mod cm-auth/go.sum ./
COPY cm-proto/go.mod cm-proto/go.sum ../cm-proto/
COPY cm-auth/.env.docker .env

RUN go mod download

COPY cm-auth/ .
COPY cm-proto/ ../cm-proto/

RUN make build

# runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates build-base

WORKDIR /app/cm-auth

COPY --from=builder /app/cm-auth/main .
COPY --from=builder /app/cm-auth/.env.docker .env

EXPOSE 8080

CMD ["./main"]

