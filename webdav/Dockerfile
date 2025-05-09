# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM golang:1.23-alpine3.21 AS builder

WORKDIR /build

COPY shared ./shared
COPY api ./api
COPY conversion ./conversion
COPY mosaic ./mosaic
COPY webdav ./webdav
COPY ui ./ui
COPY go.work go.work.sum ./

RUN go mod download

WORKDIR /build/webdav

RUN go build -o voltaserve-webdav

FROM alpine:3.21 AS runner

WORKDIR /app

COPY --from=builder /build/webdav/voltaserve-webdav ./voltaserve-webdav
COPY --from=builder /build/webdav/.env ./.env

ENTRYPOINT ["./voltaserve-webdav"]

EXPOSE 8082
