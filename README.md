# go-musthave-diploma

## Run locally

```shell
cp .env.sample .env
docker-compose --env-file .env up
```

## Nice to have

1. Change Gin logger to project zap logger.
2. Find a workaround for `Sync()` issue with zap logger (`inappropriate ioctl for device`).
3. Create a helper for error serialization instead of using raw `gin.H`.
4. Find a way to run migration while app starting process.
