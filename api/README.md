# Voltaserve API

## Getting Started

Install [Golang](https://go.dev/doc/install).

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

Run linter:

```shell
golangci-lint run
```

### Generate and Publish Documentation

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
