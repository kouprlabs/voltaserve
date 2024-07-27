# Copyright 2023 Anass Bouassaba & Daniël Sonck.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

FROM --platform=$BUILDPLATFORM rust:1.79 AS rust-builder

WORKDIR /build

ARG TARGETPLATFORM

COPY build.sh ./

RUN /build/build.sh install

COPY Cargo.toml Cargo.lock ./

# Build our dependencies first for the cache
RUN /build/build.sh deps

COPY . ./

# Now build our software
RUN /build/build.sh

FROM scratch AS runner

COPY --from=rust-builder /build/migrate /migrate

ENTRYPOINT ["/migrate"]

CMD ["up"]

EXPOSE 3000
