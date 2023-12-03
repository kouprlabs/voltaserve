# Voltaserve Identity Provider

## Getting Started

Install dependencies:

```shell
pnpm install
```

Run for development:

```shell
pnpm dev
```

Run for production:

```shell
pnpm start
```

Build Docker image:

```shell
docker build -t voltaserve/idp .
```

## Generate and Publish Documentation

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
