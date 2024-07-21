# Voltaserve WebDAV

Install [golangci-lint](https://github.com/golangci/golangci-lint).

Install [swag](https://github.com/swaggo/swag).

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

```shell
gci write -s standard -s default -s "prefix(github.com/kouprlabs)" -s "prefix(github.com/kouprlabs/voltaserve/webdav)" .
```

Format code:

```shell
gofmt -w .
```

```shell
gofumpt -w .
```

```shell
goimports -w .
```

Build Docker image:

```shell
docker build -t voltaserve/webdav .
```
