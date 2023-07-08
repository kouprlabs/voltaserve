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

Build SLE Docker image:

```shell
docker build -t voltaserve/conversion:sle -f ./Dockerfile.sle .
```

Build RHEL Docker image:

```shell
docker build -t voltaserve/conversion:rhel -f ./Dockerfile.rhel .
```
