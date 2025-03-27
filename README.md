Asset Service API
=================

Asset Service – это REST API для загрузки и скачивания данных с авторизацией. Сервис позволяет:
- Получать авторизационный токен (session-id) через эндпоинт `/api/auth`.
- Загружать данные (файлы) на сервер с помощью эндпоинта `POST /api/upload-asset/{assetName}`.
- Скачивать данные с сервера через `GET /api/asset/{assetName}`.
- Получать список всех закаченных файлов через `GET /api/assets`.
- Удалять файлы через `DELETE /api/asset/{assetName}`.
- Проверять состояние сервера через эндпоинт `/health`.

Сервис реализован на Go с использованием стандартной библиотеки и драйвера [pgx](https://github.com/jackc/pgx) для PostgreSQL. Работа сервера по протоколу HTTPS обеспечивается с помощью самоподписанных сертификатов.

------------------------------------------------------------

Структура проекта
-----------------

.
├── api/                   
│   ├── openapi.yaml       # OpenAPI спецификация API
├── certs/                 # TLS сертификаты (cert.pem, key.pem)
├── cmd/                   
│   └── main.go            # Точка входа приложения
├── internal/
│   ├── config/            # Конфигурация (чтение .env, настройки)
│   ├── db/                # Подключение к PostgreSQL через pgx
│   ├── handlers/          # HTTP-обработчики (auth, asset, health)
│   ├── models/            # Модели данных (User, Session, Asset)
│   ├── repository/        # Репозитории для работы с БД
│   └── service/           # Бизнес-логика (авторизация, валидация токенов)
├── schema.sql             # SQL-скрипт для инициализации базы данных
├── Dockerfile             # Dockerfile (multi-stage) для сборки приложения
├── docker-compose.yaml    # Docker Compose для поднятия сервиса, БД и миграций
├── .env                   # Файл с переменными окружения
└── README.md              # Этот файл

------------------------------------------------------------

Требования
----------

- Go 1.24.1 (или совместимая версия)
- PostgreSQL 15
- Docker и Docker Compose
- OpenSSL для генерации TLS-сертификатов
- (Опционально) Node.js для генерации статической документации (ReDoc CLI)

------------------------------------------------------------

Настройка
---------

### Файл .env

Создайте файл `.env` в корне проекта со следующим содержимым (значения по умолчанию):

    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=testdb

    APP_PORT=8443

    TLS_CERT_PATH=certs/cert.pem
    TLS_KEY_PATH=certs/key.pem

### Инициализация базы данных

SQL-скрипт `schema.sql` содержит схему базы данных:
- Создаются таблицы `users`, `sessions` и `assets`.
- Устанавливаются внешние ключи (ON DELETE CASCADE).
- Вставляется тестовый пользователь `alice` с паролем `secret` (хешируется с помощью MD5).

Миграция выполняется автоматически через сервис `migrate` в Docker Compose. Если база не инициализирована, можно вручную выполнить:

    docker-compose exec db psql -U postgres -d testdb -f schema.sql

### Генерация TLS сертификатов

Если сертификаты ещё не сгенерированы, выполните в терминале (например, Git Bash) из корня проекта:

    mkdir certs
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/key.pem -out certs/cert.pem -subj "/CN=localhost"

Это создаст папку `certs` с файлами `cert.pem` и `key.pem`.

------------------------------------------------------------

Запуск проекта
---------------

### С использованием Docker Compose

Убедитесь, что Docker и Docker Compose установлены. В корне проекта выполните:

    docker-compose up --build

Порядок запуска:
1. **db:** PostgreSQL запускается, данные сохраняются в volume `db_data`, healthcheck проверяет готовность.
2. **migrate:** Ждёт, пока база станет доступной, и применяет `schema.sql`.
3. **webserver:** Запускается после миграции, слушает на порту, указанном в `APP_PORT` (8443).

### Запуск без Docker

- Убедитесь, что настроены переменные окружения (или файл `.env` загружается).
- Запустите сервер командой:

  go run cmd/main.go

Приложение будет доступно по адресу `https://localhost:8443` (если настроено HTTPS) или `http://localhost:8443` (если HTTP).

------------------------------------------------------------

Использование API (с примерами cURL)
------------------------------------

> **Примечание:** Для HTTPS с самоподписанными сертификатами используйте флаг `--insecure`.

### 1. Авторизация

**Endpoint:** `POST /api/auth`

    curl -X POST -H "Content-Type: application/json" -d "{\"login\":\"alice\",\"password\":\"secret\"}" https://localhost:8443/api/auth --insecure

**Пример ответа:**

    {"token":"<ваш_токен>"}

Скопируйте полученный токен.

### 2. Загрузка данных (Upload)

**Endpoint:** `POST /api/upload-asset/{assetName}`  
Где `{assetName}` — имя файла, под которым будут сохранены данные (например, `hello`).

    curl -X POST -H "Authorization: Bearer <ваш_токен>" --data-binary "Hello, Alice!" https://localhost:8443/api/upload-asset/hello --insecure

**Пример ответа:**

    {"status":"ok"}

### 3. Скачивание данных (Download)

**Endpoint:** `GET /api/asset/{assetName}`

    curl -X GET -H "Authorization: Bearer <ваш_токен>" https://localhost:8443/api/asset/hello --insecure

**Пример ответа:**

    Hello, Alice!

### 4. Получение списка файлов

**Endpoint:** `GET /api/assets`

    curl -X GET -H "Authorization: Bearer <ваш_токен>" https://localhost:8443/api/assets --insecure

**Пример ответа:**

    {
      "assets": [
        {
          "name": "hello",
          "created_at": "2025-03-27T12:34:56Z"
        }
      ]
    }

### 5. Удаление файла

**Endpoint:** `DELETE /api/asset/{assetName}`

    curl -X DELETE -H "Authorization: Bearer <ваш_токен>" https://localhost:8443/api/asset/hello --insecure

**Пример ответа:**

    {"status":"ok"}

### 6. Healthcheck

**Endpoint:** `GET /health`

    curl -X GET https://localhost:8443/health --insecure

**Пример ответа:**

    {"status":"ok"}

------------------------------------------------------------

Документация API
----------------

Спецификация API хранится в файле [api/openapi.yaml](api/openapi.yaml).  
Вы можете открыть этот файл в [Swagger Editor](https://editor.swagger.io/) или [ReDoc](https://redocly.github.io/redoc/) для визуализации документации.

Чтобы сгенерировать статическую HTML-документацию с помощью ReDoc CLI, выполните:

    npx redoc-cli bundle api/openapi.yaml -o docs.html

Откройте файл `docs.html` в браузере для просмотра документации.

------------------------------------------------------------

Заключение
----------

Этот проект реализует:
- Авторизацию и выдачу токена.
- Загрузку, скачивание, получение списка и удаление файлов.
- Ограничение сессии до 24 часов и хранение IP адреса.
- Работа сервера по HTTPS с самоподписанными сертификатами.
- Документация API доступна через OpenAPI спецификацию.

