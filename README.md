## Сервис уведомлений (notification)

Используются:
- github.com/caarlos0/env/v9
- github.com/gin-gonic/gin
- github.com/jackc/pgx/v5
- github.com/lib/pq
- github.com/jmoiron/sqlx

### VsCode run launch
.vscode\launch.json
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

### Миграции - golang-migrate

1. Создание миграций
```
migrate create -ext sql -dir migrations -seq create_notifications_table
```

2. Выполнение миграций
Команда для применения миграций:
```
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

Команда для отката миграций:
```
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down
```


**SWAGGER**:
- github.com/swaggo/gin-swagger

Генерация yaml 
```
swag init -g cmd/notification-server/main.go --output docs  --parseDependency --parseInternal 
```