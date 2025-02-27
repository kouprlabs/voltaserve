# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY shared ./shared
COPY api ./api
COPY conversion ./conversion
COPY mosaic ./mosaic
COPY webdav ./webdav
COPY go.work go.work.sum ./

RUN go mod download

WORKDIR /build/api

RUN go build -o voltaserve-api

FROM debian:bookworm AS runner

WORKDIR /app

COPY --from=builder /build/api/voltaserve-api ./voltaserve-api
COPY --from=builder /build/api/.env ./.env
COPY --from=builder /build/api/templates ./templates

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 8080
