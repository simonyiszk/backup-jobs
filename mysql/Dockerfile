FROM golang:1.24.1-alpine AS builder

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /upload-dump

FROM alpine AS application

RUN apk add mysql-client
RUN apk add mariadb-connector-c

WORKDIR /app
COPY --from=builder /upload-dump ./upload-dump
COPY ./backup.sh ./

CMD ["/bin/sh", "./backup.sh"]
