# Copyright (c) 2023 Anass Bouassaba & Daniël Sonck.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

FROM --platform=$BUILDPLATFORM rust:1.79 AS rust-builder

WORKDIR /build

ARG TARGETPLATFORM

RUN case $TARGETPLATFORM in \
    linux/arm64) \
        dpkg --add-architecture arm64; \
        apt-get update; \
        apt-get install -y gcc-aarch64-linux-gnu musl-tools:arm64 musl-dev:arm64; \
        rustup target add aarch64-unknown-linux-musl; \
    ;; \
    linux/amd64) \
        apt-get update; \
        apt-get install -y musl-dev musl-tools; \
        rustup target add x86_64-unknown-linux-musl; \
    ;; \
    esac

COPY Cargo.toml Cargo.lock ./

# Build our dependencies first for the cache
RUN case $TARGETPLATFORM in \
    linux/arm64) \
        export CARGO_BUILD_TARGET=aarch64-unknown-linux-musl; \
        export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_MUSL_LINKER=aarch64-linux-gnu-gcc; \
    ;; \
    linux/amd64) \
        export CARGO_BUILD_TARGET=x86_64-unknown-linux-musl; \
    ;; \
    esac; \
    cargo build -p sea-orm-migration --release --locked

COPY . ./

# Now build our software
RUN case $TARGETPLATFORM in \
    linux/arm64) \
        export CARGO_BUILD_TARGET=aarch64-unknown-linux-musl; \
        export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_MUSL_LINKER=aarch64-linux-gnu-gcc; \
    ;; \
    linux/amd64) \
        export CARGO_BUILD_TARGET=x86_64-unknown-linux-musl; \
    ;; \
    esac; \
    cargo build --release --locked; \
    mv /build/target/${CARGO_BUILD_TARGET}/release/migrate /build/migrate

FROM scratch AS runner

COPY --from=rust-builder /build/migrate /migrate

ENTRYPOINT ["/migrate"]

CMD ["up"]

EXPOSE 3000
