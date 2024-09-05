# Voltaserve Console API

## Getting Started

Install dependencies:

```shell
poetry install --no-interaction --no-cache
```

Activate the virtual environment created by Poetry:

```shell
source /home/your_user/.cache/pypoetry/virtualenvs/voltaserve-console-something/bin/activate
```

Run:

```shell
poetry run python -m api.uvi --port 8086
```

Add `--reload` flag for development:

```shell
poetry run python -m api.uvi --port 8086 --reload
```

To use CockroachDB, add the `POSTGRES_PORT` environment variable:

```shell
POSTGRES_PORT=26257 poetry run python -m api.uvi --reload --port 8086
```

Lint code:

```shell
flake8 . --extend-ignore F401,W291,W503 --max-line-length 120 --extend-exclude Dockerfile
```

Format code:

```shell
black .
```
