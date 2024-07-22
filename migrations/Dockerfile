# Copyright 2023 Anass Bouassaba & Daniël Sonck.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM rust:1.79-alpine AS rust-builder

WORKDIR /build

COPY . ./

RUN apk add --no-cache musl-dev &&\
    rustc --print target-list &&\
    export ARCH=$(apk --print-arch); \
    cargo build --release --locked --target=${ARCH}-unknown-linux-musl; \
    mv /build/target/${ARCH}-unknown-linux-musl/release/migrate /build/migrate

FROM scratch AS runner

COPY --from=rust-builder /build/migrate /migrate

ENTRYPOINT ["/migrate"]

EXPOSE 3000
