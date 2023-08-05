# Voltaserve API

## Getting Started

Install [Golang](https://go.dev/doc/install), [Swag](https://github.com/swaggo/swag) and [golangci-lint](https://golangci-lint.run/usage/install).

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

### Build Docker Image

```shell
docker build -t voltaserve/api .
```

### Code Linter

```shell
golangci-lint run
```

### Generate and Publish Documentation

Format swag comments:

```shell
swag fmt
```

Generate `swagger.yml`:

```shell
swag init --output ./docs --outputTypes yaml
```

Preview (will be served at [http://127.0.0.1:5555](http://127.0.0.1:5555)):

```shell
npx @redocly/cli preview-docs --port 5555 ./docs/swagger.yaml
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ./docs/swagger.yaml --output ./docs/index.html
```
