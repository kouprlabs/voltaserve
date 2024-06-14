# Voltaserve Language

## Getting Started

Install [PDM](https://pdm-project.org).

Install Python 3.12.

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

Install spaCy model:

```shell:

spacy download xx_ent_wiki_sm
```

Run for development:

```shell
flask run --host=0.0.0.0 --port=8084 --debug
```
