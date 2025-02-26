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

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o voltaserve-mosaic

FROM alpine:3.21 AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-mosaic ./voltaserve-mosaic
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-mosaic"]

EXPOSE 8085
