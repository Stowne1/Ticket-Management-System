Ticket Management Microservice (In Progress)

• Designing and building a scalable microservice in Go for handling support ticket workflows, with a focus on RESTful API design.
• Containerized the service using Docker and orchestrated with a local Kubernetes cluster to simulate production deployment.
• Integrating AWS RDS for cloud-based database hosting and Prometheus for monitoring service health and uptime.

# Testing

## Unit/Integration Tests (Go)

To run all Go tests:

```
go test ./...
```

## API Smoke Test (End-to-End)

To run the API smoke test script (make sure your server is running on localhost:8080):

```
sh scripts/api_test.sh
```

This script will test creating, retrieving, updating, and deleting a ticket, as well as error cases.
