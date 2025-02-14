FROM golang:1.23.4 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

EXPOSE 8080

RUN go mod download
RUN go build -o bucket-app ./cmd/app/main.go