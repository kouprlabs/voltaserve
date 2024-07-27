#!/bin/bash

case $TARGETPLATFORM in
linux/arm64)
  TARGET=aarch64-unknown-linux-musl
  export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_MUSL_LINKER=aarch64-linux-gnu-gcc
  ;;
linux/amd64)
  TARGET=x86_64-unknown-linux-musl
  ;;
*)
  echo "Unknown platform $TARGETPLATFORM"
  exit 1
  ;;
esac

install_deps() {
  case $TARGETPLATFORM in
  linux/arm64)
    dpkg --add-architecture arm64
    apt-get update
    apt-get install -y gcc-aarch64-linux-gnu musl-tools:arm64 musl-dev:arm64
    rustup target add aarch64-unknown-linux-musl
  ;;
  linux/amd64)
    apt-get update
    apt-get install -y musl-dev musl-tools
    rustup target add x86_64-unknown-linux-musl
  ;;
  esac
}

build_deps() {
  cargo build -p sea-orm-migration --release --locked --target=${TARGET}
}

build_all() {
  cargo build --release --locked --target=${TARGET}
  mv /build/target/${TARGET}/release/migrate /build/migrate
}

case $1 in
install)
  install_deps
;;
deps)
  build_deps
;;
*)
  build_all
;;
esac