# Voltaserve API

## Getting Started

Install [Swag](https://github.com/swaggo/swag) and [golangci-lint](https://golangci-lint.run/usage/install).

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

Build Docker image:

```shell
docker build -t voltaserve/api .
```

Run code linter:

```shell
golangci-lint run
```

## Generate and Publish Documentation

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
