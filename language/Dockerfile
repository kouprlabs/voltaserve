FROM python:3.11-slim-bookworm

WORKDIR /app

COPY . .

RUN apt-get update

RUN apt-get install -y curl build-essential python3-poetry rust-all

RUN poetry install --no-root

RUN pip3 install torch==2.3.1

ENV PIP_DEFAULT_TIMEOUT=900

RUN poetry run spacy download xx_ent_wiki_sm

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084