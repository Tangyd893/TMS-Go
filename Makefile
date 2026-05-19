.PHONY: dev test lint typecheck build docker-up docker-down clean

dev:
	docker compose -f docker/docker-compose.yml up -d postgres redis rabbitmq minio

test:
	cd backend && go test ./... -v -cover

test-race:
	cd backend && go test -race ./... -v

lint:
	cd backend && go vet ./...
	cd frontend && npm run typecheck && npm run lint

typecheck:
	cd frontend && npm run typecheck

build:
	cd backend && go build ./cmd/server
	cd frontend && npm run build

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down

clean:
	docker compose -f docker/docker-compose.yml down -v
	rm -rf backend/tmp/ frontend/dist/
