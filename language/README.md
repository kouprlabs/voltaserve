# Voltaserve Language

## Getting Started

We assume the development environment is setup as described [here](../DEVELOPMENT.md).

Install [Pipenv](https://pipenv.pypa.io/en/latest/installation/#installing-pipenv).

Install dependencies:

```shell
pipenv install
```

Activate the environment:

```shell
pipenv shell
```

Run:

```shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```

Format code:

```shell
black .
```
