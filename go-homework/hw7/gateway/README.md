# Gateway

Gateway использует тот же HTTP API, но Ledger теперь сохраняет данные в PostgreSQL. Перед запуском должны быть доступны PostgreSQL и Redis, а миграции должны быть применены.

## Запуск

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
export REDIS_ADDR="localhost:6379"
export REDIS_DB="0"
go run .
```

Сервер работает на `http://localhost:8080`.

## Эндпойнты

- `POST /api/budgets`
- `GET /api/budgets`
- `POST /api/transactions`
- `GET /api/transactions`
- `GET /api/reports/summary?from=YYYY-MM-DD&to=YYYY-MM-DD`

Для проверки Redis вызовите отчёт два раза подряд. Первый запрос должен дать `cache miss`, второй — `cache hit` в логах приложения.
