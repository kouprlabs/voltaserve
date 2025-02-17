# Voltaserve WebDAV

Prerequisites:

- [golangci-lint v1.61.0](https://github.com/golangci/golangci-lint/releases/tag/v1.61.0)
- [gci](https://github.com/daixiang0/gci)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [swag](https://github.com/swaggo/swag)

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

Lint code:

```shell
golangci-lint run
```

Format code:

```shell
gci write -s standard -s default \
  -s "prefix(github.com/kouprlabs)" \
  -s "prefix(github.com/kouprlabs/voltaserve/webdav)" . && \
gofumpt -w . && \
gofmt -s -w . && \
golangci-lint run --fix && \
goimports -w .
```

Build Docker image:

```shell
docker build -t voltaserve/webdav .
```
