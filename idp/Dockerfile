# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM node:20-alpine

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .

RUN npm ci --legacy-peer-deps

ENTRYPOINT ["npm", "run", "start:node"]

EXPOSE 8081
