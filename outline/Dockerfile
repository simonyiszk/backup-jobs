FROM golang:1.24.1 AS builder

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backup-app

FROM alpine AS application

WORKDIR /app
COPY --from=builder /backup-app ./backup-app

CMD ["./backup-app"]
