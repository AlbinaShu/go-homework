# Go Homework 7 — PostgreSQL и Redis

В этом домашнем задании данные Ledger хранятся в PostgreSQL вместо памяти. Gateway сохраняет прежние эндпойнты бюджетов и транзакций. Отчёт `/api/reports/summary` кэшируется в Redis на 30 секунд.

## Структура

```text
hw7/
├── README.md
├── go.work
├── gateway/
│   ├── README.md
│   ├── go.mod
│   ├── main.go
│   └── internal/api/
│       ├── dto.go
│       ├── handlers.go
│       └── middleware.go
└── ledger/
    ├── README.md
    ├── go.mod
    ├── db.go
    ├── cache.go
    ├── ledger.go
    ├── report.go
    └── migrations/
        └── 00001_create_tables.sql
```

## 1. PostgreSQL

Создайте базу `cashapp`. Например, если PostgreSQL установлен локально:

```bash
createdb cashapp
```

Задайте строку подключения:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
```

Если `DATABASE_URL` не задана, используются переменные `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`. Значения по умолчанию: `localhost`, `5432`, `postgres`, `postgres`, `cashapp`.

## 2. Goose и миграции

Установка Goose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Из каталога `hw7` применить миграцию:

```bash
goose -dir ./ledger/migrations postgres "$DATABASE_URL" up
```

Откатить последнюю миграцию:

```bash
goose -dir ./ledger/migrations postgres "$DATABASE_URL" down
```

Миграция создаёт таблицы `budgets` и `expenses`.

## 3. Redis

Самый простой запуск через Docker:

```bash
docker run -p 6379:6379 --name cashcraft-redis -d redis:7-alpine
```

Переменные Redis:

```bash
export REDIS_ADDR="localhost:6379"
export REDIS_DB="0"
export REDIS_PASSWORD=""
```

`REDIS_ADDR` и `REDIS_DB` можно не задавать: по умолчанию используются `localhost:6379` и `0`.

## 4. Зависимости и запуск

Первый раз загрузите зависимости:

```bash
cd ledger
go mod tidy
cd ../gateway
go mod tidy
cd ..
```

Запуск Gateway из `hw7/gateway`:

```bash
cd gateway
go run .
```

При успешном старте появятся сообщения о подключении к PostgreSQL и Redis, затем сервер будет слушать порт `8080`.

## 5. Проверка API

Создать бюджет:

```bash
curl -i -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"Food","limit":5000}'
```

Получить бюджеты:

```bash
curl -i http://localhost:8080/api/budgets
```

Добавить транзакцию в пределах бюджета:

```bash
curl -i -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":1000,"category":"Food","description":"Groceries","date":"2026-11-02"}'
```

Попытка превысить бюджет должна вернуть `409 Conflict`:

```bash
curl -i -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":4500,"category":"Food","description":"Restaurant","date":"2026-11-02"}'
```

Получить транзакции:

```bash
curl -i http://localhost:8080/api/transactions
```

## 6. Проверка Redis-кэша отчёта

Выполните один и тот же запрос два раза подряд:

```bash
curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
curl "http://localhost:8080/api/reports/summary?from=2026-11-01&to=2026-11-30"
```

При первом запросе в логах будет `report cache miss`, результат считается через PostgreSQL и записывается в Redis. При втором запросе в течение 30 секунд будет `report cache hit`.

## 7. Проверка сохранения данных

Остановите Gateway (`Ctrl+C`) и снова выполните:

```bash
cd gateway
go run .
```

После перезапуска запросы `GET /api/budgets` и `GET /api/transactions` должны вернуть ранее добавленные данные, потому что они находятся в PostgreSQL.

## Что реализовано

- PostgreSQL хранит бюджеты и расходы.
- `SetBudget` использует UPSERT.
- `AddTransaction` проверяет сумму расходов перед вставкой и возвращает `budget exceeded` при превышении лимита.
- `ListBudgets` и `ListTransactions` читают данные из PostgreSQL.
- Пул БД: максимум 10 открытых и 5 idle-соединений.
- Redis настраивается через `REDIS_ADDR`, `REDIS_DB`, `REDIS_PASSWORD`.
- Отчёт по категориям кэшируется в Redis на 30 секунд.
