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

COPY . .

RUN apk update
RUN apk add --no-cache build-base

RUN pip3 install pipx
ENV PATH="/root/.local/bin:$PATH"

RUN pipx install pdm --python $(which python)

RUN pdm install --prod --no-editable
RUN .venv/bin/python3 -m ensurepip

ENTRYPOINT ["pdm", "run", "flask", "run", "--host=0.0.0.0", "--port=8084"]

EXPOSE 8084