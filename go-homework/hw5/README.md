# Go Homework 5

В домашнем задании реализован HTTP API в сервисе Gateway. Gateway использует модуль Ledger из того же проекта. Данные хранятся в памяти.

## Структура

```text
hw5/
├── README.md
├── gateway/
│   ├── go.mod
│   ├── main.go
│   └── internal/api/
│       ├── dto.go
│       ├── handlers.go
│       └── middleware.go
└── ledger/
    ├── go.mod
    └── ledger.go
```

## Запуск

```bash
cd gateway
go run .
```

Сервер запускается на `http://localhost:8080`.

## API

### Создать или обновить бюджет

```bash
curl -X POST http://localhost:8080/api/budgets \
  -H "Content-Type: application/json" \
  -d '{"category":"Food","limit":5000}'
```

### Получить бюджеты

```bash
curl http://localhost:8080/api/budgets
```

### Создать транзакцию

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount":1000,"category":"Food","description":"Groceries","date":"2026-10-31"}'
```

### Получить транзакции

```bash
curl http://localhost:8080/api/transactions
```

Если транзакция превышает установленный бюджет, сервер возвращает статус `409 Conflict` и JSON:

```json
{"error":"budget exceeded"}
```

Ошибки валидации возвращаются со статусом `400 Bad Request`.
