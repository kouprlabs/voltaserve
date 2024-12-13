# Voltaserve Identity Provider

## Getting Started

Run for development:

```shell
deno task dev
```

Run for production:

```shell
deno task start
```

Lint code:

```shell
deno lint .
```

Format code:

```shell
deno fmt .
```

Build Docker image:

```shell
docker build -t voltaserve/idp .
```

Preview (will be served at [http://localhost:7777](http://localhost:7777)):

```shell
deno -A npm:@redocly/cli preview-docs --port 7777 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
deno -A npm:@redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
