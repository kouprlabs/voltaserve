# Voltaserve API

## Getting Started

We assume the development environment is setup as described [here](../DEVELOPMENT.md).

### Build and Run

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

### Setup the Linter

Install [golangci-lint](https://golangci-lint.run).

Run the following and make sure there are no issues:

```shell
golangci-lint run
```

### Generate and Publish Documentation

Install [swag](https://github.com/swaggo/swag):

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Format swag comments:

```shell
swag fmt
```

We suppose that the [api-docs](https://github.com/voltaserve/api-docs) repository is cloned locally at: `../../api-docs/`.

Generate `swagger.yml`:

```shell
swag init --output ../../api-docs/ --outputTypes yaml
```

Preview (will be served at [http://127.0.0.1:5555](http://127.0.0.1:5555)):

```shell
npx @redocly/cli preview-docs --port 5555 ../../api-docs/swagger.yaml
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ../../api-docs/swagger.yaml --output ../../api-docs/index.html
```

Now you can open a PR in the [api-docs](https://github.com/voltaserve/api-docs) repository with your current changes.
