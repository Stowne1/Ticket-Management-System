# Stage 1: build the React frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: build the Go binary (with frontend/dist embedded)
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built frontend into place so the Go embed picks it up
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
WORKDIR /app/main
RUN go build -o /go/bin/ticket-app

# Stage 3: minimal runtime image
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /go/bin/ticket-app ./
EXPOSE 8080
CMD ["./ticket-app"]
