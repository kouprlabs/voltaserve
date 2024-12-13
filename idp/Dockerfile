# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM denoland/deno:debian-2.1.4 AS builder

WORKDIR /build

COPY src ./src
COPY templates ./templates
COPY .env .
COPY deno.json .
COPY deno.lock .

RUN deno task compile

FROM debian:bookworm-slim AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-idp ./voltaserve-idp
COPY --from=builder /build/.env ./.env
COPY --from=builder /build/templates ./templates

ENTRYPOINT ["./voltaserve-idp"]

EXPOSE 8081
