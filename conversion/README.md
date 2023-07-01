# Voltaserve Conversion

## Getting Started

Install [Golang](https://go.dev/doc/install).

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

### Docker Images

Build SLE 15 Docker image:

```shell
docker build -t voltaserve/conversion -f ./Dockerfile.sle15 .
```

Build RHEL 9 Docker image:

```shell
docker build -t voltaserve/conversion -f ./Dockerfile.rhel9 .
```
