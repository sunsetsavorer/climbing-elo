DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/climbing_elo?sslmode=disable

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

.PHONY: test
test:
	go test ./...

%:
	@:
