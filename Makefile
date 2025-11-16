APP_NAME := pull-review
PORT := 8080
DB_DSN ?= postgres://postgres:postgres@localhost:5432/pull_review?sslmode=disable

.PHONY: build run lint test migrate-up migrate-down compose-up compose-down

build: 
	go build -o bin/${APP_NAME} ./cmd/server

run: 
	go run ./cmd/server

lint:
	golangci-lint run

test:
	go test ./...

migrate-up:
	docker run --rm -v $$PWD/migrations:/migrations --network host \
	migrate/migrate -path=/migrations -database "$(DB_DSN)" up

migrate-down:
	docker run --rm -v $$PWD/migrations:/migrations --network host \
	migrate/migrate -path=/migrations -database "$(DB_DSN)" down 1

compose-up:
	docker compose -f deploy/docker-compose.yml --env-file configs/.env.example up -d

compose-down:
	docker compose -f deploy/docker-compose.yml down -v