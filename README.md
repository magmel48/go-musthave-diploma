# go-musthave-diploma-tpl

## Run locally

```shell
cp .env.sample .env
docker-compose --env-file .env up
```

## Known bugs

After only first `docker-compose up` run schema changes will not be applied, re-run.
