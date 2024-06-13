FROM python:3.11-slim-bookworm

WORKDIR /app

COPY . .

RUN apt-get update

RUN apt-get install -y curl build-essential python3-poetry rust-all

RUN poetry install --no-root

RUN pip3 install torch==2.3.1

ENV PIP_DEFAULT_TIMEOUT=900

RUN poetry run spacy download xx_ent_wiki_sm
RUN poetry run spacy download zh_core_web_trf
RUN poetry run spacy download de_core_news_lg
RUN poetry run spacy download en_core_web_trf
RUN poetry run spacy download fr_core_news_lg
RUN poetry run spacy download pt_core_news_lg
RUN poetry run spacy download es_core_news_lg

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084