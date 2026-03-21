.PHONY: setup migrate-% up up-deps down-deps up-normal up-testing down

ENV_FILE := .env.dev

setup: migrate-reset migrate-up

migrate-%:
	go run . --env-file $(ENV_FILE) migrate $* $(ARGS)

migrate-up:
migrate-down:
migrate-reset:
migrate-create:

seed:
	go run . --env-file $(ENV_FILE) seed

up:
	air -c .air.toml

up-deps:
	docker compose up postgres

down-deps:
	docker compose down

up-normal:
	docker compose up postgres backend -d

up-testing:
	docker compose up postgres backend-testing -d

down:
	docker compose down
