# Voltaserve Language

## Getting Started

Install dependencies:

```shell
poetry install
```

Run:

```shell
poetry run flask run --host=0.0.0.0 --port=8084
```

Add `--debug` flag for development:

```shell
poetry run flask run --host=0.0.0.0 --port=8084 --debug
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
