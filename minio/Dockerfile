FROM golang:1.24.1-alpine AS builder

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY *.go ./

RUN apk add -U --no-cache ca-certificates
RUN go install github.com/minio/mc@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o /upload-dump

FROM alpine AS application

WORKDIR /app
COPY --from=builder /go/bin/mc  /usr/bin/mc
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /upload-dump ./upload-dump
COPY ./backup.sh ./

CMD ["/bin/sh", "./backup.sh"]
