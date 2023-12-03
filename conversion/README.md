# Voltaserve Conversion

Install [golangci-lint](https://golangci-lint.run/usage/install).

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

Run code linter:

```shell
golangci-lint run
```

Build Docker image:

```shell
docker build -t voltaserve/conversion .
```
