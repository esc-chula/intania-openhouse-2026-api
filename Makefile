.PHONY: setup migrate-% up up-deps down-deps docker-build docker-run

ENV_FILE := .env.dev

setup: migrate-reset migrate-up

migrate-%:
	go run . --env-file $(ENV_FILE) migrate $* $(ARGS)

migrate-up:
migrate-down:
migrate-reset:
migrate-create:

up:
	air -c .air.toml

up-deps:
	docker compose --env-file $(ENV_FILE) up

down-deps:
	docker compose --env-file $(ENV_FILE) down

docker-build:
	docker build -t intania-openhouse-2026-api -f Dockerfile .

docker-run: docker-build
	docker run --env-file $(ENV_FILE) intania-openhouse-2026-api
