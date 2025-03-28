# Voltaserve Identity Provider

## Getting Started

Run:

```shell
deno task start
```

Check code:

```shell
deno check .
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

Preview (will be served at [http://localhost:19091](http://localhost:19091)):

```shell
deno -A npm:@redocly/cli preview-docs --port 19091 ./docs/swagger.json
```

Generate the final static HTML documentation:

```shell
deno -A npm:@redocly/cli build-docs ./docs/swagger.json --output ./docs/index.html
```
