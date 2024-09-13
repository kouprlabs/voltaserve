# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM oven/bun:1-alpine

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY bun.lockb .
COPY tsconfig.json .

RUN bun install --frozen-lockfile --production

ENTRYPOINT ["bun", "run", "start"]

EXPOSE 8081
