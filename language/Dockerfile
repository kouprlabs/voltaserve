# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM python:3.11.9-bookworm

WORKDIR /app

COPY . .

RUN apt-get update &&\
    apt-get install -y curl build-essential curl rust-all

RUN pip install poetry &&\
    poetry config virtualenvs.create false

RUN poetry install --no-root

ENTRYPOINT ["flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084