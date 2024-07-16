# balance service

## Configuration
Configuration only via environment
- DB_HOST - Database host (Default localhost)
- DB_PORT - Database port (Default 5432)
- DB_USER - Database user
- DB_DATABASE - Database name
- DB_PASSWORD - Database password
- APP_HOST - Host for web server
- APP_PORT - Port for web server
- APP_ENV - Env (default dev)

## Running

```
docker compose up
```
For run migrations in docker set `RUN_MIGRATIONS=1`

### Run migrations plain

```
go run ./cmd/migrator
```


## Endpoints
Open path "/docs" for try SwaggerUI
