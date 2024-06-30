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

Lint code:

```shell
bun run lint
```

Format code:

```shell
bun run format
```

Build Docker image:

```shell
docker build -t voltaserve/idp .
```

## Generate Documentation

Generate `swagger.json`:

```shell
npm run swagger-autogen && mv ./swagger.json ./docs
```

Preview (will be served at [http://localhost:7777](http://localhost:7777)):

```shell
npx @redocly/cli preview-docs --port 7777 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
