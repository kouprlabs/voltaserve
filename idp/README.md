# Voltaserve Identity Provider

## Getting Started

Install dependencies:

```shell
bun i
```

Run for development:

```shell
bun run dev
```

Run for production:

```shell
bun run start
```

Build Docker image:

```shell
docker build -t voltaserve/idp .
```

## Generate and Publish Documentation

Generate `swagger.json`:

```shell
pnpm swagger-autogen && mv ./swagger.json ./docs
```

Preview (will be served at [http://127.0.0.1:7777](http://127.0.0.1:7777)):

```shell
npx @redocly/cli preview-docs --port 7777 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
