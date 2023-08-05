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

### Download spaCy

On SLE / openSUSE Leap:

```shell
python3.11 -m spacy download xx_ent_wiki_sm
```

On RHEL:

```shell
python -m spacy download xx_ent_wiki_sm
```

On other systems, including macOS and Windows:

```shell
python3 -m spacy download xx_ent_wiki_sm
```

### Run for Development

```shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```

### Build Docker Image

```shell
docker build -t voltaserve/language .
```

### Code Formatter

```shell
black .
```
