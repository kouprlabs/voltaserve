# Voltaserve UI

## Getting Started

Install [Node.js 18.x](https://nodejs.org).

Enable [pnpm](https://pnpm.io):

```shell
corepack enable
```

Install dependencies:

```shell
pnpm install
```

Run for development:

```shell
pnpm run dev
```

### Docker Images

Build SLE 15 Docker image:

```shell
docker build -t voltaserve/ui -f ./Dockerfile.sle15 .
```

Build RHEL 9 Docker image:

```shell
docker build -t voltaserve/ui -f ./Dockerfile.rhel9 .
```
