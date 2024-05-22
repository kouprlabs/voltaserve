FROM python:3.12-bookworm

WORKDIR /app

COPY . .

RUN curl -sSL https://install.python-poetry.org | python3 -
ENV PATH="/root/.local/bin:${PATH}"

RUN poetry install

RUN poetry run spacy download xx_ent_wiki_sm
RUN poetry run spacy download zh_core_web_trf
RUN poetry run spacy download de_dep_news_trf
RUN poetry run spacy download en_core_web_trf
RUN poetry run spacy download fr_dep_news_trf
RUN poetry run spacy download it_core_news_lg
RUN poetry run spacy download ja_core_news_trf
RUN poetry run spacy download nl_core_news_lg
RUN poetry run spacy download pt_core_news_lg
RUN poetry run spacy download ru_core_news_lg
RUN poetry run spacy download es_dep_news_trf
RUN poetry run spacy download sv_core_news_lg

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084