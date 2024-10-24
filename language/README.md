# Voltaserve Language

## Getting Started

Use Python 3.12:

```shell
pdm use cpython@3.12.3
```

Install dependencies:

```shell
pdm install
```

Activate the virtual environment created by PDM:

```shell
source .venv/bin/activate
```

Make sure PIP is available:

```shell
python3 -m ensurepip
```

Run:
```shell
flask run --host=0.0.0.0 --port=8084
```

Add `--debug` flag for development:

```shell
flask run --host=0.0.0.0 --port=8084 --debug
```

Lint code:

```shell
flake8 .
```

Format code:
```shell
black .
```