# Voltaserve Language

## Getting Started

Install [Python](https://www.python.org) and [Pipenv](https://pipenv.pypa.io).

Install dependencies:

```shell
pipenv install
```

```shell
python3 -m spacy download xx_ent_wiki_sm
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
