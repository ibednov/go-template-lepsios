# go-template-lepsios

Skeleton Go API. Clone / copy into a new service, rename module and `APP_NAME`.

## Stack

- Go 1.27, Gin, uber/fx, cobra (`server` | `migrate` | `job`)
- goose, sqlboiler, envconfig via go-lepsios
- zerolog; optional Vector → VictoriaLogs via `make logs`

## Layout

```
cmd/                 server | migrate | job
internal/
  app/               fx graph, engine, routes
  config/
  domain/<feature>/  deep vertical slice
  sidecar/           postgres, redis
configs/
migrations/
deployments/{local,prod,res}/
```

Domain slice:

```
internal/domain/<feature>/
  entity/
  errors/
  service/
  db/postgres/repository/
  db/postgres/convert/
  transport/api/http/v1/handler|request|response/
  fx.go
```

## Local

```bash
cp go.work.example go.work   # expects sibling ../go-lepsios
cp configs/.env.example configs/.env
make surround               # postgres+redis + migrate
make run                    # API on :8080
make logs                   # optional VictoriaLogs+Vector
```

`make surround` does **not** start logs. Logs are opt-in via separate compose.

## Replace placeholders

- `github.com/ibednov/go-template-lepsios` → your module
- `APP_NAME` / compose project / ports in Makefile and `.env`
- Delete `internal/domain/example` after first real domain exists


## Files + migrations

`go-lepsios/files` is always wired (`FILES_DRIVER=local|s3|minio`).

| Env | Default | Meaning |
|---|---|---|
| `MIGRATION_RUN` | `true` | Entrypoint runs `migrate up` before `server` |
| `MIGRATION_BACKUP` | `true` | If pending migrations — `pg_dump` into files storage |
| `MIGRATION_BACKUP_PREFIX` | `backups/pre-migrate` | Path prefix inside FILES storage |

Backup uses the same `files.Adapter` as the app (local dir or S3/MinIO). Fail closed: backup error aborts migrate.

## Env families

| Prefix | Role |
|---|---|
| `PORT` / `HTTP_HOST_PORT` | listen vs host publish |
| `SERVER_*` | HTTP timeouts |
| `LOG_*` | level, vector URL, local UI/vector ports |
| `DB_*` | Postgres |
| `REDIS_*` | Redis |
| `FILES_*` | go-lepsios files.Adapter |
| `MIGRATION_*` | run on start + pre-migrate backup |
| `APP_*` / `*_CONTAINER_NAME` | local compose names |

