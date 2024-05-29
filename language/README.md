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
poetry run spacy download xx_ent_wiki_sm
poetry run spacy download zh_core_web_trf
poetry run spacy download de_core_news_lg
poetry run spacy download en_core_web_trf
poetry run spacy download fr_core_news_lg
poetry run spacy download it_core_news_lg
poetry run spacy download ja_core_news_trf
poetry run spacy download nl_core_news_lg
poetry run spacy download pt_core_news_lg
poetry run spacy download ru_core_news_lg
poetry run spacy download es_core_news_lg
poetry run spacy download sv_core_news_lg
```

On Apple Silicon or Intel Macs with supported AMD GPUs, do the following to install a hardware accelerated version of PyTorch:

https://developer.apple.com/metal/pytorch

```shell
pip3 install --pre torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/nightly/cpu
```

Run for development:

```shell
poetry run flask run --host=0.0.0.0 --port=8084 --debug
```
