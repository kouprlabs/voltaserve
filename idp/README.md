# Voltaserve Identity Provider

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

Build SLE Docker image:

```shell
docker build -t voltaserve/idp:sle -f ./Dockerfile.sle .
```

Build RHEL Docker image:

```shell
docker build -t voltaserve/idp:rhel -f ./Dockerfile.rhel .
```

### Generate and Publish Documentation

Generate `swagger.json`:

```shell
pnpm run swagger-autogen && mv ./swagger.json ./docs
```

Preview (will be served at [http://127.0.0.1:7777](http://127.0.0.1:7777)):

```shell
npx @redocly/cli preview-docs --port 7777 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
