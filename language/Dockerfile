FROM python:3.12-bookworm

WORKDIR /app

COPY . .

RUN curl -sSL https://install.python-poetry.org | python3 -
ENV PATH="/root/.local/bin:${PATH}"

RUN poetry install

RUN poetry run spacy download xx_ent_wiki_sm

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084