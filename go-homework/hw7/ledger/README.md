# Ledger — PostgreSQL и Redis

Ledger больше не хранит бюджеты и транзакции в памяти. Данные читаются и записываются в PostgreSQL, а отчёт по категориям кэшируется в Redis на 30 секунд.

## PostgreSQL

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/cashapp?sslmode=disable"
goose -dir ./migrations postgres "$DATABASE_URL" up
```

Откат:

```bash
goose -dir ./migrations postgres "$DATABASE_URL" down
```

Если `DATABASE_URL` не задана, используются `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`.

## Redis

```bash
export REDIS_ADDR="localhost:6379"
export REDIS_DB="0"
export REDIS_PASSWORD=""
```

Ключ отчёта имеет вид `report:summary:<from>:<to>`. TTL — 30 секунд.
