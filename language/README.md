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

Download spaCy model on SLE 15:

```shell
python3.11 -m spacy download xx_ent_wiki_sm
```

Download spaCy model on RHEL 9:

```shell
python3 -m spacy download xx_ent_wiki_sm
```

Run for development:

```shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```

### Docker Images

Build SLE 15 Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.sle15 .
```

Build RHEL 9 Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.rhel9 .
```

### Code Formatter

```shell
black .
```
