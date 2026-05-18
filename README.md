# TMS-Go

TMS-Go is a transportation management system skeleton based on Golang, Vue 3, PostgreSQL, and Docker-managed middleware.

## Project Structure

```text
TMS-Go/
  backend/     Golang microservices workspace
  frontend/    Vue 3 frontend application
  docker/      Docker Compose, Nginx, environment examples
  docs/        Design documents
  testing/     Test and smoke-check scripts
  todo/        Local todo workspace, ignored by Git
```

## Local Middleware

```powershell
cd docker
docker compose up -d postgres redis rabbitmq minio
```

## Backend Gateway Seed

```powershell
cd backend
go mod tidy
go run ./cmd/server
```

## Frontend

```powershell
cd frontend
npm install
npm run dev
```

## Health Check

```text
GET http://localhost:8080/health
GET http://localhost:8080/ready
```
