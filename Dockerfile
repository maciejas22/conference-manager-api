FROM golang:1.22-bookworm as builder

WORKDIR /usr/src/app

COPY Makefile ./
COPY .env .env.prod ./

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN GO_ENV=prod make build

FROM debian:bookworm

WORKDIR /usr/src/app
RUN apt-get update -y && apt-get install -y ca-certificates 


COPY --from=builder /usr/src/app /usr/src/app

ENV GO_ENV=prod
CMD ["./bin/conference-manager-api"]
