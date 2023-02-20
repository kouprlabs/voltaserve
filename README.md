<!-- markdownlint-disable MD033 MD041 -->
<p align="center">
  <img height="70" src="assets/brand.svg"/>
  <h1 align="center">Voltaserve</h1>
</p>

## Getting Started

Install [Docker](https://docs.docker.com/get-docker) and [Docker Compose](https://docs.docker.com/compose/install).

### Run for Development

Start infrastructure services:

```sh
docker compose up
```

Follow instructions for running:

- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)

### Run for Production

Optionally build Docker images locally, if not they will be downloaded from [Docker Hub](https://hub.docker.com) in the second step:

```sh
docker-compose -f ./docker-compose.prod.yml build
```

```sh
docker-compose -f ./docker-compose.prod.yml up
```

## Licensing

Voltaserve is released under the [The MIT License](./LICENSE).
