# go-musthave-diploma

## Run locally

```shell
cp .env.sample .env
docker-compose --env-file .env up
```

## Known bugs

1. After only first `docker-compose up` run schema changes will not be applied, re-run.

## Nice to have

1. Change Gin logger to project zap logger.
2. Fix known bugs.
3. Find a workaround for `Sync()` issue with zap logger (`inappropriate ioctl for device`).
