# Voltaserve Shared

## Getting Started

Prerequisites:

- [golangci-lint](https://github.com/golangci/golangci-lint/releases/tag/v1.61.0) (v1.61.0)
- [gci](https://github.com/daixiang0/gci)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [gotestsum](https://github.com/gotestyourself/gotestsum)
- [swag](https://github.com/swaggo/swag)

Lint code:

```shell
golangci-lint run
```

Format code:

```shell
gci write -s standard -s default \
  -s "prefix(github.com/kouprlabs)" \
  -s "prefix(github.com/kouprlabs/voltaserve/shared)" . && \
gofumpt -w . && \
gofmt -s -w . && \
goimports -w . && \
golangci-lint run --fix && \
swag fmt
```