FROM python:3.11-slim-bookworm

WORKDIR /app

COPY . .

RUN apt-get update

RUN apt-get install -y curl build-essential

RUN apt-get update && \
  apt-get install -y curl && \
  curl https://sh.rustup.rs -sSf | sh -s -- -y && \
  /root/.cargo/bin/rustup default stable

RUN curl -sSL https://install.python-poetry.org | python3 -

ENV PATH="/root/.cargo/bin:/root/.local/bin:${PATH}"

RUN poetry install --no-root

RUN poetry run spacy download xx_ent_wiki_sm
RUN poetry run spacy download zh_core_web_trf
RUN poetry run spacy download de_core_news_lg
RUN poetry run spacy download en_core_web_trf
RUN poetry run spacy download fr_core_news_lg
RUN poetry run spacy download it_core_news_lg
RUN poetry run spacy download ja_core_news_trf
RUN poetry run spacy download nl_core_news_lg
RUN poetry run spacy download pt_core_news_lg
RUN poetry run spacy download ru_core_news_lg
RUN poetry run spacy download es_core_news_lg
RUN poetry run spacy download sv_core_news_lg

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084