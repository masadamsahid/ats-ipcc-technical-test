# Running Migration

Firstly, you have to install [`golang-migrate`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) in your system.

# Commands

## To generate migration files:

```bash
migrate create -ext=sql -seq -dir="./db/migrations" FILE_NAME
```

## To run migrations

```bash
migrate -path="./db/migrations" -database="postgres://postgres:postgres@localhost:5432/tsilodot_db?sslmode=disable" up
```

## To rollback

```bash
migrate -path="./db/migrations" -database="postgres://postgres:postgres@localhost:5432/tsilodot_db?sslmode=disable" down
```

Rollback to specific version (e.g. to 1st migration)

```bash
migrate -path="./db/migrations" -database="postgres://postgres:postgres@localhost:5432/tsilodot_db?sslmode=disable" down 1
```

For more other command, visit [`golang-migrate cli docs`](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).