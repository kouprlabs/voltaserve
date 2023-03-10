# Development

## Generate and Publish Documentation

We suppose that the [idp-docs](https://github.com/voltaserve/idp-docs) repository is cloned locally at: `../../idp-docs/`.

Generate `swagger.json`:

```sh
npm run swagger-autogen && mv ./swagger.json ../../idp-docs
```

Preview (will be served at [http://127.0.0.1:7777](http://127.0.0.1:7777)):

```sh
npx @redocly/cli preview-docs --port 7777 ../../idp-docs/swagger.json
```

Generate the final static HTML documentation:

```sh
npx @redocly/cli build-docs ../../idp-docs/swagger.json --output ../../idp-docs/index.html
```

Now you can open a PR in the [idp-docs](https://github.com/voltaserve/idp-docs) repository with your current changes.
