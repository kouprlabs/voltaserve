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

Lint code:

```shell
poetry run flake8 .
```

Format code:

```shell
poetry run black .
```

Sort imports:

```shell
poetry run isort .
```
