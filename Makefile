include .env
export

export PROJECT_ROOT=$(shell pwd)

access-rights:
	sudo chmod -R 777 ./out/pgdata/18/docker
	sudo chmod -R 777 ./internal/migrations

webchat-run:
	@go run cmd/main.go

env-up:
	@docker compose up -d webchat-postgres
env-down:
	@docker compose stop webchat-postgres
env-cleanup:
	@read -p "Are you sure you want to cleanup the environment? (y/n): " ans;\
	if [ "$$ans" = "y" ]; then\
		docker compose down webchat-postgres && \
		sudo rm -rf out/pgdata && \
		echo "Environment cleaned up successfully";\
	else \
		echo "Environment cleanup cancelled"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq={seq}"; \
		exit 1; \
	fi; \
	docker compose run --rm webchat-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"
migrate-up:
	@make migrate-action action=up
migrate-down:
	@make migrate-action action=down
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action={action}"; \
		exit 1; \
	fi; \

	docker compose run --rm webchat-postgres-migrate \
       	-path /migrations \
       	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@webchat-postgres:5432/${POSTGRES_DB}?sslmode=disable \
        "$(action)"
