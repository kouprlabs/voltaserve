FROM python:3.12-alpine

WORKDIR /app

ADD pyproject.toml .

RUN apk update &&\
    pip install --upgrade pip &&\
    apk add --no-cache --virtual .build-deps gcc build-base libffi-dev libretls-dev cargo &&\
    pip install poetry

COPY . .

RUN poetry config virtualenvs.create false &&\
    poetry install --no-interaction --no-ansi &&\
    apk del .build-deps &&\
    rm -rf /root/.cache /root/.cargo

ENTRYPOINT ["python","-m", "api.uvi"]

