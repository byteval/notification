.PHONY: test test-unit test-integration test-coverage build run clean lint swagger migrate-up migrate-down

BINARY_NAME=notification-server
MAIN_FILE=cmd/notification-server/main.go

# Тесты
test:
	go test -v ./...

test-unit:
	go test -v -tags=unit ./...

test-integration:
	go test -v -tags=integration ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии: coverage.html"

test-short:
	go test -v -short ./...

# Сборка
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux $(MAIN_FILE)

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe $(MAIN_FILE)

# Запуск
run:
	go run $(MAIN_FILE)

run-dev:
	gin -i -b bin/gin-bin -a 8080 run $(MAIN_FILE)

# Очистка
clean:
	go clean
	rm -rf bin/
	rm -f coverage.out coverage.html

# Линтер
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

# Swagger
swagger:
	swag init -g $(MAIN_FILE) --output docs --parseDependency --parseInternal

# Миграции
migrate-up:
	migrate -path migrations -database "postgres://php2go:php2go@localhost:5432/notifications?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://php2go:php2go@localhost:5432/notifications?sslmode=disable" down

# Docker
docker-build:
	docker build -f Dockerfile.dev -t $(BINARY_NAME):dev .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Зависимости
deps:
	go mod download
	go mod tidy

deps-update:
	go get -u ./...
	go mod tidy

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  test              - Запустить все тесты"
	@echo "  test-unit         - Запустить только unit тесты"
	@echo "  test-integration  - Запустить только интеграционные тесты"
	@echo "  test-coverage     - Запустить тесты с отчетом о покрытии"
	@echo "  test-short        - Запустить короткие тесты"
	@echo "  build             - Собрать приложение"
	@echo "  run               - Запустить приложение"
	@echo "  run-dev           - Запустить с hot reload (требует gin)"
	@echo "  clean             - Очистить артефакты сборки"
	@echo "  lint              - Запустить линтер"
	@echo "  lint-fix          - Запустить линтер с автоисправлением"
	@echo "  swagger           - Сгенерировать Swagger документацию"
	@echo "  migrate-up        - Применить миграции"
	@echo "  migrate-down      - Откатить миграции"
	@echo "  docker-build      - Собрать Docker образ"
	@echo "  docker-run        - Запустить Docker контейнеры"
	@echo "  docker-stop       - Остановить Docker контейнеры"
	@echo "  deps              - Установить зависимости"
	@echo "  deps-update       - Обновить зависимости" 