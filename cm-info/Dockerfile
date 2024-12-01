# builder 
FROM golang:latest AS builder

WORKDIR /app/cm-info

COPY cm-info/go.mod cm-info/go.sum ./
COPY cm-proto/go.mod cm-proto/go.sum ../cm-proto/
COPY cm-info/.env.docker .env

RUN go mod download

COPY cm-info/ .
COPY cm-proto/ ../cm-proto/

RUN make build

# runner
FROM alpine:latest

RUN apk --no-cache add ca-certificates build-base

WORKDIR /app/cm-info

COPY --from=builder /app/cm-info/main .
COPY --from=builder /app/cm-info/.env.docker .env

EXPOSE 8080

CMD ["./main"]

