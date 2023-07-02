# Voltaserve WebDAV

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

Run for production:

```shell
pnpm run start
```

### Docker Images

Build SLE / openSUSE Leap Docker image:

```shell
docker build -t voltaserve/webdav -f ./Dockerfile.sle .
```

Build RHEL Docker image:

```shell
docker build -t voltaserve/webdav -f ./Dockerfile.rhel .
```
