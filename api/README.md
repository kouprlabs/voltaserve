# Voltaserve API

## Getting Started

Prerequisites:

- [golangci-lint](https://github.com/golangci/golangci-lint/releases/tag/v1.61.0) (v1.61.0)
- [gci](https://github.com/daixiang0/gci)
- [gofumpt](https://github.com/mvdan/gofumpt)
- [gotestsum](https://github.com/gotestyourself/gotestsum)
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
  -s "prefix(github.com/kouprlabs/voltaserve/api)" . && \
gofumpt -w . && \
gofmt -s -w . && \
goimports -w . && \
golangci-lint run --fix && \
swag fmt
```

Run tests:

```shell
gotestsum --format testdox --packages="./..."
```

Build Docker image:

```shell
docker build -t voltaserve/api .
```

## Generate Documentation

Generate `swagger.yml`:

```shell
swag init --parseDependency --output ./docs --outputTypes yaml
```

Preview (will be served at [http://localhost:19090](http://localhost:19090)):

```shell
bunx @redocly/cli preview-docs --port 19090 ./docs/swagger.yaml
```

Generate the final static HTML documentation:

```shell
bunx @redocly/cli build-docs ./docs/swagger.yaml --output ./docs/index.html
```
