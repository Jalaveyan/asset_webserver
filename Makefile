.PHONY: build run migrate docker-build docker-run docker-stop

# Локальная сборка (без Docker)
build:
	go build -o bin/webserver cmd/webserver/main.go

run:
	./bin/webserver

# Применение SQL-схемы
migrate:
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f schema.sql

# Сборка Docker-образа
docker-build:
	docker build -t go-asset-service:latest .

# Локальный запуск через docker-compose
docker-run:
	docker-compose up --build

docker-stop:
	docker-compose down
