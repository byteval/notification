## Сервис уведомлений (notification)

Основные функции:

1. REST API - управление шаблонами и уведомлениями.

2. Отправка сообщений на электронную почту (SMTP).

3. Обработка входящей почты (IMAP).

4. Возврат обработанных ответов исходным клиентам по вебхуку.

### Структура проекта
```
notification/
├── cmd/                                # Точка входа
├── config/                             # Конфигурация
├── docs/                               # Swagger
├── migrations/                         # Миграции
├── internal/                           
│    ├── app/                           # Операции / Use Cases
│    │   └── notification/              # Контекст - Уведомления
│    │       └── create/                # Создание уведомления
│    │           ├── handler.go         # Хэндлер операции
│    │           ├── mapper.go          # Маппинг ToDomain/ToResponse
│    │           └── dto.go             # DTO Request/Response
│    ├── container/            
│    │   └── container.go               # Контейнер DI
│    ├── domain/                        # Domain: Сущности, Интерфейсы сервисов, Бизнес-правила
│    │   └── notification/              # Контекст - Уведомления
│    │         ├── notification.go      # Модель (Notification)
│    │         └── ports/               # Интерфейсы сервисов
│    │             └── repository.go    # Интерфейс репозитория Notification
│    ├── transport/
│    │   └── http/                      # API 
│    │       ├── handlers/              # Обработчики HTTP запросов
│    │       │   └── notifications/
│    │       │       ├── create.go      # POST /notifications
│    │       │       └── get.go         # GET /notifications/{id}    
│    │       ├── middleware/
│    │       │   ├── middleware.go   
│    │       │   └── validator.go   
│    │       └── router.go
│    └── infrastructure/                # Реализации (БД, SMTP и т.д.)
│        └── postgres/
│            ├── database.go            # Подключение к СУБД POSTGRES
│            └── notification_repository.go
└── test/ 
```
[Полное описание структуры проекта](./structure.md)

### Используются:
- github.com/caarlos0/env/v9
- github.com/gin-gonic/gin
- github.com/jackc/pgx/v5
- github.com/lib/pq
- github.com/jmoiron/sqlx


### Миграции - golang-migrate

1. Создание миграций
```
migrate create -ext sql -dir migrations -seq create_notifications_table
```

2. Выполнение миграций
Команда для применения миграций:
```
migrate -path migrations -database "postgres://php2go:php2go@localhost:5432/notifications?sslmode=disable" up
```

Команда для отката миграций:
```
migrate -path migrations -database "postgres://php2go:php2go@localhost:5432/notifications?sslmode=disable" down
```


### Swagger
- github.com/swaggo/gin-swagger

Генерация yaml 
```
swag init -g cmd/notification-server/main.go --output docs  --parseDependency --parseInternal 
```

http://localhost:8080/swagger/index.html


### Linter

https://golangci-lint.run/

```
golangci-lint run 
```

### VsCode run launch
.vscode\launch.json с поддержкой `.env`
```
{
	"version": "0.2.0",
	"configurations": [
		{
			"name": "Launch",
			"type": "go",
			"request": "launch",
			"mode": "auto",
			"envFile": "${workspaceFolder}/.env",
			"program": "${workspaceFolder}/cmd/notification-server/main.go"
		}
	]
}
```