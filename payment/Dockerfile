FROM python:3.12-slim AS builder

ARG APP_VERSION

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

WORKDIR /app
COPY poetry.lock pyproject.toml /app/

RUN pip install poetry && \
    poetry config virtualenvs.create false && \
    poetry version "$APP_VERSION" && \
    poetry install --no-interaction --no-ansi

COPY /src /app/

FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

WORKDIR /app

COPY --from=builder /usr/local/lib/python3.12/site-packages/ /usr/local/lib/python3.12/site-packages/
COPY --from=builder /usr/local/bin/uvicorn /usr/local/bin/uvicorn

COPY --from=builder /app /app

EXPOSE 8080
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]