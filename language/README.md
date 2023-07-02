# Voltaserve Language

## Getting Started

Install [Python](https://www.python.org) and [Pipenv](https://pipenv.pypa.io).

On RHEL, run the following:

```shell
pipenv --python /usr/bin/python
```

Install dependencies:

```shell
pipenv install
```

Activate environment:

```shell
pipenv shell
```

Download spaCy model on SLE / openSUSE Leap:

```shell
python3.11 -m spacy download xx_ent_wiki_sm
```

Download spaCy model on RHEL:

```shell
python -m spacy download xx_ent_wiki_sm
```

Run for development:

```shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```

### Docker Images

Build SLE / openSUSE Leap Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.sle .
```

Build RHEL Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.rhel .
```

### Code Formatter

```shell
black .
```
