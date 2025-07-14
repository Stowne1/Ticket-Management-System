# Start with the official Golang image for building
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/main
RUN go build -o /go/bin/ticket-app

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/bin/ticket-app ./
EXPOSE 8080
ENV PORT=8080
CMD ["./ticket-app"]