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

RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then \
  .venv/bin/pip3 install https://download.pytorch.org/whl/cpu/torch-2.3.1+cpu-cp312-cp312-linux_x86_64.whl; \
  elif [ "$BUILDPLATFORM" = "linux/arm64" ]; then \
  .venv/bin/pip3 install https://download.pytorch.org/whl/cpu/torch-2.3.1-cp312-cp312-manylinux_2_17_aarch64.manylinux2014_aarch64.whl; \
  fi

RUN .venv/bin/python3 -m spacy download xx_ent_wiki_sm
RUN .venv/bin/python3 -m spacy download zh_core_web_trf
RUN .venv/bin/python3 -m spacy download de_core_news_lg
RUN .venv/bin/python3 -m spacy download en_core_web_trf
RUN .venv/bin/python3 -m spacy download fr_core_news_lg
RUN .venv/bin/python3 -m spacy download it_core_news_lg
RUN .venv/bin/python3 -m spacy download ja_core_news_trf
RUN .venv/bin/python3 -m spacy download nl_core_news_lg
RUN .venv/bin/python3 -m spacy download pt_core_news_lg
RUN .venv/bin/python3 -m spacy download ru_core_news_lg
RUN .venv/bin/python3 -m spacy download es_core_news_lg
RUN .venv/bin/python3 -m spacy download sv_core_news_lg

ENTRYPOINT ["pdm", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084