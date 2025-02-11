ENV_FILE = .env

ENV_VARS = \
    POSTGRES_DB=avito \
    POSTGRES_USER=user \
    POSTGRES_PASSWORD=password \
    POSTGRES_HOST=db \
    POSTGRES_PORT=5432 \
	SSL_MODE=disable \

env:
	@$(eval SHELL:=/bin/bash)
	@printf "%s\n" $(ENV_VARS) > $(ENV_FILE)
	@echo "$(ENV_FILE) file created"

run:
	@chmod +x docker/scripts/entrypoint.sh
	@docker compose up --build -d

runl:
	@chmod +x docker/scripts/entrypoint.sh
	@docker compose up --build

off:
	@docker compose down

db:
	@docker compose up --build -d db

logs:
	@docker compose logs