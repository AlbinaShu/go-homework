# Go Homework 2

ДЗ состоит из двух сервисов:

- **Gateway** — HTTP-сервис
- **Ledger** — сервис бизнес-логики для работы с транзакциями

## Требования

- Go версии 1.22+

---

## Gateway

Gateway — простой HTTP-сервис, который отвечает на GET-запрос `/ping`.

### Запуск

Перейти в папку сервиса:

```bash
cd gateway
```

Запустить:

```bash
go run .
```

После запуска сервис будет доступен на порту 8080.

### Проверка работы

В отдельном терминале выполнить:

```bash
curl http://localhost:8080/ping
```

Ожидаемый ответ:

```text
pong
```

---

## Ledger

Ledger — сервис для работы с финансовыми транзакциями.

Реализовано:

- структура `Transaction`;
- хранение транзакций в памяти через `[]Transaction`;
- добавление транзакций через функцию `AddTransaction`;
- получение списка транзакций через функцию `ListTransactions`;
- автоматическое создание ID транзакции.

### Запуск

Перейти в папку сервиса:

```bash
cd ledger
```

Запустить:

```bash
go run .
```

Ожидаемый вывод:

```text
Ledger service started
Transactions:
{1 150.5 Food Dinner 2026-09-27}
{2 50 Transport Taxi 2026-09-27}
{3 1000 Shopping Shoes 2026-09-27}
```

---

## Структура проекта

```text
hw2/
├── gateway/
│   ├── go.mod
│   └── main.go
├── ledger/
│   ├── go.mod
│   └── main.go
└── README.md
```

## Проверка проекта

### Gateway

```bash
cd gateway
go run .
```

В отдельном терминале:

```bash
curl http://localhost:8080/ping
```

Ответ:

```text
pong
```

### Ledger

```bash
cd ledger
go run .
```

Сервис должен вывести сообщение о запуске и список тестовых транзакций.