ENV_LOCAL=./envs/local.env
# Загружаем переменные из .env-файла
include $(ENV_LOCAL)
export $(shell sed 's/=.*//' $(ENV_LOCAL))


# Docker команды
docker-up:
	docker-compose --file $(COMPOSE_FILE) up -d

docker-down:
	docker-compose --file $(COMPOSE_FILE) down

docker-logs:
	docker-compose --file $(COMPOSE_FILE) logs -f

docker-restart:
	docker-compose --file $(COMPOSE_FILE) down && docker-compose --file $(COMPOSE_FILE) up -d

# Работа с БД и миграциями
new-migration:
	goose -dir db/migrations create $(name) sql

migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir db/migrations up

migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir db/migrations down

migrate-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DATABASE_URL)" goose -dir db/migrations status
