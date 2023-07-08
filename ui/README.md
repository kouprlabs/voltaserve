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

Build SLE Docker image:

```shell
docker build -t voltaserve/ui:sle -f ./Dockerfile.sle .
```

Build RHEL Docker image:

```shell
docker build -t voltaserve/ui:rhel -f ./Dockerfile.rhel .
```
