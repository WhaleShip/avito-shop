ENV_FILE = .env

ENV_VARS = \
    POSTGRES_DB=avito \
    POSTGRES_USER=user \
    POSTGRES_PASSWORD=password \
    POSTGRES_HOST=db \
    POSTGRES_PORT=5432 \
	SSL_MODE=disable \
	JWTSECRET=dontHackMePls \
	PGBOUNCER_HOST=pgbouncer \
	PGBOUNCER_PORT=6432 \

env:
	@$(eval SHELL:=/bin/bash)
	@printf "%s\n" $(ENV_VARS) > $(ENV_FILE)
	@echo "$(ENV_FILE) file created"

run:
	@docker compose up

runl:
	@go run cmd/app/main.go

off:
	@docker compose down

build:
	@docker compose build

db:
	@docker compose up --build -d db pgbouncer

logs:
	@docker compose logs

lint:
	@golangci-lint run

cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

cover-html:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

test:
	@go test ./...

test-int:
	@go test -v -tags=integration ./tests/integration_tests