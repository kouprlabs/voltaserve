# Voltaserve Conversion

## Getting Started

Install [Golang](https://go.dev/doc/install) and [golangci-lint](https://golangci-lint.run/usage/install).

### Build and Run

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

### Code Linter

```shell
golangci-lint run
```

### Build Docker Image

```shell
docker build -t voltaserve/conversion .
```
