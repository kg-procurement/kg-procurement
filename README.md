# KG Procurement

---

## Getting Started

Download Go version `1.23.1` and run `make tidy` then `go mod vendor`

For config, copy-paste `config.example.jsonc` to `config.jsonc`

## Running

Run db: `make docker-up`

## PgAdmin

After pgadmin container is running, open `localhost:15432`, then create server.

Set the following on the `Connection` tab:

- `host`: `postgres-local` (the database service name on compose)
- `port`: `5432`
- `username`: `<POSTGRES_USER>`
- `password`: `<POSTGRES_PASSWORD>`

## Add schema migrations

- Create migrations: `goose create <migration_name> sql`
- Move migration to migrations directory (Manually): `mv <"FILE_NAME">.sql migrations/<"FILE_NAME">.sql`

## Execute migrations

- Apply all available migrations: `make migrate-up`
- Role back single migrations from the current version: `make migrate-down`
