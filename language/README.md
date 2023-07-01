# Voltaserve Language

## Getting Started

Install [Python](https://www.python.org) and [Pipenv](https://pipenv.pypa.io).

### Install Dependencies

```shell
pipenv install
```

### Activate the Environment

```shell
pipenv shell
```

### Download spaCy Model

On openSUSE:

```shell
python3.11 -m spacy download xx_ent_wiki_sm
```

On macOS:

```shell
python3 -m spacy download xx_ent_wiki_sm
```

### Run

```shell
FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug
```

### Format Code

```shell
black .
```

### Build Docker Image

Build SLE 15 Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.sle15 .
```

Build RHEL 9 Docker image:

```shell
docker build -t voltaserve/language -f ./Dockerfile.rhel9 .
```
