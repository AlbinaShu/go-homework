# Go Homework 6

В домашнем задании добавлены unit-тесты для Ledger и HTTP-интеграционные тесты для Gateway. Данные по-прежнему хранятся в памяти.

## Структура

```text
hw6/
├── README.md
├── gateway/
│   ├── go.mod
│   ├── main.go
│   └── internal/api/
│       ├── dto.go
│       ├── handlers.go
│       ├── handlers_test.go
│       └── middleware.go
└── ledger/
    ├── go.mod
    ├── ledger.go
    └── ledger_test.go
```

## Запуск Gateway

```bash
cd gateway
go run .
```

Сервер запускается на `http://localhost:8080`.

## Тестирование

Запуск unit-тестов Ledger:

```bash
cd ledger
go test ./... -cover
```

Запуск тестов Gateway:

```bash
cd gateway
go test ./... -cover
```

Создание отчёта о покрытии:

```bash
go test ./... -cover -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html
```

Покрытие тестами:

- Ledger — 83.3%
- Gateway — 57.7% (общее покрытие модуля)
- Gateway internal/api — 76.3%