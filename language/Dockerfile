# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM python:3.12-alpine

WORKDIR /app

RUN apk update &&\
    apk add --no-cache build-base &&\
    pip3 install pipx

COPY . .

ENV PATH="/root/.local/bin:$PATH"

RUN pipx install pdm --python $(which python) &&\
    pdm install --prod --no-editable &&\
    .venv/bin/python3 -m ensurepip

ENTRYPOINT ["pdm", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084