# builder 
FROM golang:latest AS builder

WORKDIR /app/cm-conferences

COPY cm-conferences/go.mod cm-conferences/go.sum ./
COPY cm-proto/go.mod cm-proto/go.sum ../cm-proto/
COPY cm-conferences/.env.docker .env

RUN go mod download

COPY cm-conferences/ .
COPY cm-proto/ ../cm-proto/

RUN make build

# runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates build-base

WORKDIR /app/cm-conferences

COPY --from=builder /app/cm-conferences/main .
COPY --from=builder /app/cm-conferences/.env.docker .env

EXPOSE 8080

CMD ["./main"]

