# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM node:22.14-alpine AS ui-builder

WORKDIR /build

COPY src ./src
COPY public ./public
COPY index.html .
COPY postcss.config.cjs .
COPY tailwind.config.cjs .
COPY vite.config.ts .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY tsconfig.node.json .

RUN npm ci
RUN npm run build

FROM golang:1.23-alpine3.21 AS go-builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY --from=ui-builder /build/dist ./dist
COPY config ./config
COPY main.go .

RUN go build -o voltaserve-ui

FROM alpine:3.21 AS runner

WORKDIR /app

COPY --from=go-builder /build/voltaserve-ui ./voltaserve-ui
COPY .env .

ENTRYPOINT ["./voltaserve-ui"]

EXPOSE 3000
