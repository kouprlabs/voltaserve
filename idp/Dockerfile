# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM denoland/deno:alpine-2.1.4

WORKDIR /app

COPY src ./src
COPY templates ./templates
COPY .env .
COPY deno.json .
COPY deno.lock .

ENTRYPOINT ["deno", "task", "start"]

EXPOSE 8081
