# Voltaserve Language

## Getting Started

Install [Poetry](https://python-poetry.org).

Install dependencies:

```shell
poetry install
```

Spawn a shell within the project's virtual environment:

```shell
poetry shell
```

Install spaCy model:

```shell:
python3 -m spacy download xx_ent_wiki_sm
```

Run for development:

```shell
flask run --host=0.0.0.0 --port=8084 --debug
```
