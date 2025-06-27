**Структура каталогов**:
```
notification/
├── cmd/
│   └── notification-server/
│       └── main.go                     # Инициализация конфига, логгера, DI-контейнера
├── config/                             # Конфигурация
│   └── config.go                       # Структура конфига проекта (параметры базы данных, SMTP, IMAP, Websocket, Webhook)
├── docs/                               # Swagger
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── migrations/                         # Миграции
├── internal/
│    ├── app/                           # Операции / Use Cases
│    │   ├── notification/              # Уведомления
│    │   │   ├── create/                # Создание уведомления
│    │   │   │   ├── handler.go         # Создание уведомления
│    │   │   │   ├── notifier.go        # Отправка уведомления
│    │   │   │   ├── mapper.go          # Маппинг ToDomain/ToResponse
│    │   │   │   ├── template_renderer.go # Рендер уведомления
│    │   │   │   └── dto.go             # DTO Request/Response
│    │   │   ├── get/                   # Получение уведомления
│    │   │   │   ├── handler.go          
│    │   │   │   ├── mapper.go          
│    │   │   │   └── dto.go    
│    │   │   └── mailprocessor/         # Получение ответов на уведомления
│    │   │       ├── handler.go          
│    │   │       ├── validator.go          
│    │   │       └── notifier.go             
│    │   └── layout/                    # Шаблоны уведомлений
│    │       ├── common/                # Общие DTO и функции
│    │       │   ├── mapper.go
│    │       │   └── dto.go             # Общие DTO
│    │       ├── create/
│    │       │   ├── handler.go
│    │       │   ├── validator.go
│    │       │   └── dto.go
│    │       ├── update/
│    │       │   ├── handler.go
│    │       │   └── mapper.go
│    │       ├── get/
│    │       │   ├── handler.go     
│    │       │   └── dto.go
│    │       ├── list/
│    │       │   ├── handler.go
│    │       │   ├── mapper.go
│    │       │   └── dto.go
│    │       └── delete/
│    │           ├── dto.go
│    │           └── handler.go
│    ├── container/            
│    │   └── container.go               # Контейнер DI
│    ├── domain/                        # Domain: Сущности, Интерфейсы сервисов, Бизнес-правила
│    │   ├── notification/              
│    │   │     ├── errors.go            # Структура с описанием ошибок
│    │   │     ├── notification.go      # Модель (Notification)
│    │   │     └── ports/               # Интерфейсы сервисов
│    │   │         ├── repository.go    # Интерфейс репозитория Notification
│    │   │         ├── smtp.go          # Отправка уведомления SMTP
│    │   │         ├── webhook.go       # Отправка уведомления webhook
│    │   │         └── websocket.go     # Отправка уведомления websocket
│    │   └── layout/                    # Шаблоны
│    │       ├── layout.go              # Модель (Layout)
│    │       ├── errors.go
│    │       └── ports/
│    │           └── repository.go
│    ├── transport/
│    │   └── http/                      # API 
│    │       ├── handlers/              # Обработчики HTTP запросов
│    │       │   ├── healthcheck.go
│    │       │   ├── notifications/
│    │       │   │   ├── create.go      # POST /notifications
│    │       │   │   └── get.go         # GET /notifications/{id}    
│    │       │   └── layouts/           
│    │       │       ├── create.go      # POST /layouts
│    │       │       ├── update.go      # PUT /layouts/{id}
│    │       │       ├── get.go         # GET /layouts/{id}
│    │       │       ├── list.go        # GET /layouts
│    │       │       └── delete.go      # DELETE /layouts/{id}
│    │       ├── middleware/
│    │       │   ├── middleware.go   
│    │       │   └── validator.go   
│    │       └── router.go
│    ├── infrastructure/                # Реализации (БД, SMTP и т.д.)
│    │   ├── postgres/
│    │   │   ├── database.go            # Подключение к СУБД POSTGRES
│    │   │   ├── notification_repository.go
│    │   │   └── layout_repository.go   
│    │   └── notifiers/                 # Реализации сервисов
│    │       ├── smtp/                  # Сервис отправки почты по SMTP
│    │       │   └── notifier.go
│    │       └── webhook/               # Сервис отправки webhook
│    │           └── notifier.go
├── pkg/                                # Дополнительные пакеты
│    └── logger 
├── test/ 
```
