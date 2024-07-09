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

Build Docker image:

```shell
docker build -t voltaserve/webdav .
```
