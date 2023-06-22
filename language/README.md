# Voltaserve Language

## Getting Started

Install [Python](https://www.python.org) and [Pipenv](https://pipenv.pypa.io).

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
