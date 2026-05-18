# Backend Microservices

Golang microservices workspace for TMS-Go.

## Layout

```text
cmd/              Gateway seed entrypoint for early bootstrap
configs/          Shared configuration examples
internal/         Gateway seed internals
services/         Business microservice directories
shared/           Cross-service technical libraries
migrations/       Migration examples and service-level migrations
tests/            Backend integration tests
```

The initial `cmd/server` is a lightweight gateway seed with health endpoints. As services are implemented, each service should live under `services/{service-name}` with its own `cmd`, `internal`, `configs`, `migrations`, and tests.

## Run

```powershell
go mod tidy
go run ./cmd/server
```
