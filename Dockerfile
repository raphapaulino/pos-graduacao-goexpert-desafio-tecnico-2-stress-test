FROM golang:1.22.3 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stresstestapp

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stresstestapp .
ENTRYPOINT ["./stresstestapp"]