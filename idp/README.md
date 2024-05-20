# Voltaserve Identity Provider

## Getting Started

Install dependencies:

```shell
npm i --legacy-peer-deps
```

Run for development:

```shell
npm run dev
```

Run for production:

```shell
npm run start
```

Build Docker image:

```shell
docker build -t voltaserve/idp .
```

## Generate and Publish Documentation

Generate `swagger.json`:

```shell
npm run swagger-autogen && mv ./swagger.json ./docs
```

Preview (will be served at [http://localhost:7777](http://localhost:7777)):

```shell
npmx @redocly/cli preview-docs --port 7777 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
npmx @redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
