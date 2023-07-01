# Voltaserve Identity Provider

## Getting Started

Install [Node.js 18.x](https://nodejs.org).

Enable [pnpm](https://pnpm.io):

```shell
corepack enable
```

### Build and Run

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

Build SLE 15 Docker image:

```shell
docker build -t voltaserve/idp -f ./Dockerfile.sle15 .
```

Build RHEL 9 Docker image:

```shell
docker build -t voltaserve/idp -f ./Dockerfile.rhel9 .
```

### Generate and Publish Documentation

We suppose that the [idp-docs](https://github.com/voltaserve/idp-docs) repository is cloned locally at: `../../idp-docs/`.

Generate `swagger.json`:

```shell
npm run swagger-autogen && mv ./swagger.json ../../idp-docs
```

Preview (will be served at [http://127.0.0.1:7777](http://127.0.0.1:7777)):

```shell
npx @redocly/cli preview-docs --port 7777 ../../idp-docs/swagger.json
```

Generate the final static HTML documentation:

```shell
npx @redocly/cli build-docs ../../idp-docs/swagger.json --output ../../idp-docs/index.html
```

Now you can open a PR in the [idp-docs](https://github.com/voltaserve/idp-docs) repository with your current changes.
