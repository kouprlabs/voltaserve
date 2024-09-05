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
poetry run python -m api.uvi
```

Lint code:

```shell
flake8 . --extend-ignore F401,W291 --max-line-length 120 --extend-exclude Dockerfile
```

Format code:

```shell
black .
```
