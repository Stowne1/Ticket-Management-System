# Ticket Management System

A full-stack ticket management service: a Go REST API with JWT authentication, a React + TypeScript frontend embedded in the Go binary, and Postgres for storage. Ships as a single Docker image and runs on AWS.

**Live demo:** http://54.184.140.247:8080

## Features

- **Go + Gin REST API** for creating, listing, updating, and deleting tickets
- **JWT authentication** with register and login; ticket routes are protected
- **React + TypeScript frontend** built with Vite and embedded into the Go binary, so one container serves both the UI and the API
- **Bun ORM with versioned SQL migrations** that run automatically on startup
- **Automated tests:** 30 Go unit/integration tests plus an end-to-end API smoke test
- **CI with GitHub Actions:** every push to `main` runs the Go tests, builds the app, and runs the smoke test against a real Postgres service
- **Multi-stage Docker build** (Node → Go → Alpine) for a small runtime image

## Architecture

```
Browser ──► EC2 (t4g.micro, ARM) ──► RDS PostgreSQL
              └─ Docker container     (private, no public access)
                 pulled from ECR
```

- **ECR** stores the Docker image
- **EC2** runs the container; an IAM role lets it pull from ECR without stored credentials
- **RDS PostgreSQL** holds the data. It isn't publicly accessible, and its security group only accepts connections from the app server's security group
- Secrets (`POSTGRES_DSN`, `JWT_SECRET`) are passed as environment variables at runtime and never committed

## API

| Method | Route | Auth |
|---|---|---|
| GET | `/api/health` | No |
| POST | `/api/register` | No |
| POST | `/api/login` | No |
| POST | `/api/tickets` | JWT |
| GET | `/api/tickets` | JWT |
| GET | `/api/tickets/:id` | JWT |
| PUT | `/api/tickets/:id` | JWT |
| DELETE | `/api/tickets/:id` | JWT |

## Run it locally

**Prerequisites:** Docker Desktop and Git. Go 1.24+ and Node 20+ are only needed to run without Docker.

```bash
git clone https://github.com/Stowne1/Ticket-Management-System.git
cd Ticket-Management-System
docker compose up --build
```

Then open http://localhost:8080. Docker Compose starts Postgres, waits until it's healthy, and starts the app, which runs migrations automatically.

### Environment variables

| Variable | Purpose |
|---|---|
| `POSTGRES_DSN` | Postgres connection string, e.g. `postgres://user:pass@host:5432/ticketdb?sslmode=disable` (use `sslmode=require` for RDS) |
| `JWT_SECRET` | Secret for signing tokens. Use a long random value in production (`openssl rand -hex 32`) |

The values in `docker-compose.yml` are for local development only.

## Tests

Go unit and integration tests:

```bash
go test ./...
```

End-to-end API smoke test (with the server running on `localhost:8080`):

```bash
sh scripts/api_test.sh
```

## Database migrations

Migration files live in `migrations/` and are applied automatically when the app starts. To run them manually:

```bash
export POSTGRES_DSN="postgres://ticketuser:ticketpass@localhost:5432/ticketdb?sslmode=disable"
go run ./cmd/migrate db init
go run ./cmd/migrate db migrate
```

## Deploying to AWS

```bash
# Build and push the image to ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com
docker build -t ticket-app .
docker tag ticket-app:latest <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com/ticket-app:latest
docker push <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com/ticket-app:latest

# On the EC2 instance
sudo docker run -d --name ticket-app --restart unless-stopped -p 8080:8080 \
  -e POSTGRES_DSN='postgres://ticketuser:<PASSWORD>@<RDS_ENDPOINT>:5432/ticketdb?sslmode=require' \
  -e JWT_SECRET='<RANDOM_SECRET>' \
  <ACCOUNT_ID>.dkr.ecr.us-west-2.amazonaws.com/ticket-app:latest
```

## Roadmap

- CI auto-deploy to EC2 on merge to `main`
- HTTPS on a custom domain
- Ticket assignment and user roles
- OpenAPI/Swagger documentation
- Prometheus metrics and monitoring
