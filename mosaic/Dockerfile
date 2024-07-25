# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN go mod download

COPY . .

RUN go build -o voltaserve-mosaic

FROM golang:1.22-alpine AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-mosaic ./voltaserve-mosaic
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-mosaic"]

EXPOSE 8085
