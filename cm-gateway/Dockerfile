# builder 
FROM golang:latest AS builder

WORKDIR /app/cm-gateway

COPY cm-gateway/go.mod cm-gateway/go.sum ./
COPY cm-proto/go.mod cm-proto/go.sum ../cm-proto/
COPY cm-gateway/.env.docker .env

RUN go mod download

COPY cm-gateway/ .
COPY cm-proto/ ../cm-proto/

RUN make build

# runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates build-base

WORKDIR /app/cm-gateway

COPY --from=builder /app/cm-gateway/main .
COPY --from=builder /app/cm-gateway/.env.docker .env

EXPOSE 8080

CMD ["./main"]

