# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM python:3.12-slim-bookworm

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

ENTRYPOINT ["poetry", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084