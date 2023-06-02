# Development

## Linter

Install [golangci-lint](https://golangci-lint.run).

Run the following and make sure there are no issues:

```shell
golangci-lint run
```

Format swag comments:

```shell
swag fmt
```

## Generate and Publish Documentation

Install [swag](https://github.com/swaggo/swag):

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

We suppose that the [api-docs](https://github.com/voltaserve/api-docs) repository is cloned locally at: `../../api-docs/`.

Generate `swagger.yml`:

```shell
swag init --output ../../api-docs/ --outputTypes yaml
```

Preview (will be served at [http://127.0.0.1:5555](http://127.0.0.1:5555)):

```shell
npx @redocly/cli preview-docs --port 5555 ../../api-docs/swagger.yaml
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ../../api-docs/swagger.yaml --output ../../api-docs/index.html
```

Now you can open a PR in the [api-docs](https://github.com/voltaserve/api-docs) repository with your current changes.
