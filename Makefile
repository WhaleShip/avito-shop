ENV_FILE = .env

ENV_VARS = \
    POSTGRES_DB=avito \
    POSTGRES_USER=user \
    POSTGRES_PASSWORD=password \
    POSTGRES_HOST=db \
    POSTGRES_PORT=5432 \
	SSL_MODE=disable \
	JWTSECRET=dontHackMePls \

env:
	@$(eval SHELL:=/bin/bash)
	@printf "%s\n" $(ENV_VARS) > $(ENV_FILE)
	@echo "$(ENV_FILE) file created"

run:
	@docker compose up --build

runl:
	@go run cmd/app/main.go

off:
	@docker compose down

build:
	@docker compose build

db:
	@docker compose up --build -d db

logs:
	@docker compose logs

lint:
	@golangci-lint run
