# Voltaserve Conversion

Prerequisites:
- [golangci-lint](https://github.com/golangci/golangci-lint)
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
  -s "prefix(github.com/kouprlabs/voltaserve/conversion)" . && \
gofumpt -w . && \
gofmt -s -w . && \
goimports -w .
```

Build Docker image:

```shell
docker build -t voltaserve/conversion .
```

## Generate Documentation

Generate `swagger.yml`:

```shell
swag init --output ./docs --outputTypes yaml
```

Preview (will be served at [http://localhost:19093](http://localhost:19093)):

```shell
bunx @redocly/cli preview-docs --port 19093 ./docs/swagger.yaml
```

Generate the final static HTML documentation:

```shell
bunx @redocly/cli build-docs ./docs/swagger.yaml --output ./docs/index.html
```
