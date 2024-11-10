# Voltaserve Console API

## Getting Started

Install dependencies:

```shell
poetry install --no-interaction --no-cache
```

Run:

```shell
poetry run python -m api.uvi
```

Add `--reload` flag for development:

```shell
poetry run python -m api.uvi --reload
```

To use CockroachDB, add the `POSTGRES_PORT` environment variable:

```shell
POSTGRES_PORT=26257 poetry run python -m api.uvi --reload
```

Lint code:

```shell
flake8 .
```

Format code:

```shell
black .
```
