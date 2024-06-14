FROM python:3.12-alpine

WORKDIR /app

COPY . .

RUN apk update
RUN apk add --no-cache build-base

RUN pip3 install pipx
ENV PATH="/root/.local/bin:$PATH"

RUN pipx install pdm --python $(which python)

RUN pdm install --prod --no-editable
RUN .venv/bin/python3 -m ensurepip
RUN .venv/bin/python3 -m spacy download xx_ent_wiki_sm

ENTRYPOINT ["pdm", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084