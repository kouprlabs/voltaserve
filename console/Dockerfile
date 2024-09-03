# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM python:3.12-alpine

ARG USERNAME=voltaserve
ARG GROUPNAME=$USERNAME
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN addgroup \
    --gid "$USER_UID" \
    "$GROUPNAME" &&\
    adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$USER_UID" \
    $USERNAME

WORKDIR /app

ADD pyproject.toml .

RUN apk update &&\
    pip install --upgrade pip &&\
    apk add --no-cache --virtual .build-deps gcc build-base libffi-dev libretls-dev cargo &&\
    pip install poetry &&\
    poetry config virtualenvs.create false &&\
    poetry config installer.max-workers 8 &&\
    poetry config installer.parallel true

COPY . .

RUN poetry install --no-interaction --no-ansi --no-cache &&\
    apk del .build-deps &&\
    rm -rf /root/.cache /root/.cargo &&\
    chown $USERNAME:$GROUPNAME /app

USER $USERNAME

ENTRYPOINT ["python", "-m", "api.uvi"]

EXPOSE 8086