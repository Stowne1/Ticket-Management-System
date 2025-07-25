
# Ticket Management Microservice 


Ticket Management Microservice is a lightweight backend service built with Go for managing support ticket workflows. It’s designed for developers and teams who need a scalable, production-ready foundation for handling tickets, user management, and API integrations. With built-in database migrations, Dockerized Postgres, automated tests, and CI/CD readiness, it’s ideal for learning modern backend practices or powering real-world applications.

## Features
- Go + Gin: Fast, idiomatic REST API for ticket management.
- Bun ORM & Migrations: Versioned, repeatable database schema using Bun and SQL migration files.
- Postgres via Docker Compose: Local development and testing with Dockerized Postgres.
- API Smoke Tests: Automated end-to-end tests for all ticket endpoints.
- CI/CD Ready: Easily integrates with GitHub Actions or other CI tools.
## Getting Started
1. Prerequisites
- Go
- Docker Desktop
- Git

2. Clone the Repository
git clone <your-repo-url>
cd Ticket-Management-System-1

3. Start Services
docker compose up -d

4. Run Database Migrations
POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"

5. Run the App
POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"

Or let Docker Compose run it for you (default).

## Running Tests

To run tests, run the following command

```bash
  Testing
Unit/Integration Tests (Go)
To run all Go tests:
go test ./...

API Smoke Test (End-to-End)
To run the API smoke test script (make sure your server is running on localhost:8080):
sh scripts/api_test.sh
```

This script will test creating, retrieving, updating, and deleting a ticket, as well as error cases.
## Database Migrations

- Migrations are managed with Bun.
- Migration files are in the migrations/ directory.

Create a New Migration
POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"

Apply Migrations
POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"

Rollback
POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"

## In Progress / TODO

- User authentication and authorization
- User management and ticket assignment
- API documentation (OpenAPI/Swagger)
- Prometheus metrics and monitoring
- Cloud deployment