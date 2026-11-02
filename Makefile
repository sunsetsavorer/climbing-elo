-include configs/.env.local

DATABASE_DSN = $(DB_DRIVER)://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: install-tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_DSN)" up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_DSN)" down 1

.PHONY: test
test:
	go test ./...

.PHONY: run-bot
run-bot:
	go run ./cmd/bot

.PHONY: build
build:
	go build -o build/bot ./cmd/bot

.PHONY: compose-up
compose-up:
	docker compose up --build -d

.PHONY: compose-logs
compose-logs:
	docker compose logs -f

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: compose-reset
compose-reset:
	docker compose down -v

%:
	@:
