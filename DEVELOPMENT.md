# Voltaserve Development

## Prerequisites

Install [Deno](https://deno.com).

Install [Bun](https://bun.sh).

Install [Node.js](https://nodejs.org).

Install [Go](https://go.dev).

Install [Rust](https://www.rust-lang.org).

Install [Poetry](https://python-poetry.org/docs/#installation).

Install [Python](https://www.python.org) 3.11 with [pyenv](https://github.com/pyenv/pyenv).

## Run Infrastructure Services

### Using Docker

```shell
docker compose up -d \
    postgres \
    minio \
    meilisearch \
    redis \
    maildev
```

Run the [migrations/migrate]() tool in the newly created database.

### From Binaries

#### CockroachDB

Download the [binary archive](https://www.cockroachlabs.com/docs/releases) and extract the archive.

Start CockroachDB:

```shell
./cockroach start-single-node --insecure --http-addr=0.0.0.0:18080
```

Using DBeaver or any PostgreSQL GUI, connect with `root` and no password, then create a user and database:

```sql
CREATE DATABASE voltaserve;
CREATE USER voltaserve;
GRANT ALL PRIVILEGES ON DATABASE voltaserve TO voltaserve;
```

Run the [migrations/migrate]() tool in the newly created database.

#### Redis

Download the [source archive](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-from-source/) and follow the guide to build from the source.

Start Redis:

```shell
./src/redis-server
```

#### MinIO

Download the [binary](https://min.io/docs/minio/macos/index.html) and move it into its own folder like: `minio`.

Assign execute permission to the binary:

```shell
chmod +x ./minio
```

Start MinIO:

```shell
MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve ./minio server ./data --console-address ":9001"
```

#### Meilisearch

Download the [binary](https://github.com/meilisearch/meilisearch/releases/tag/v1.8.3) and move it into its own folder like: `meilisearch`.

This [guide](https://www.meilisearch.com/docs/learn/getting_started/installation) might be useful.

Assign execute permission to the binary:

```shell
chmod +x ./meilisearch
```

_Note: the binary will have a suffix matching the architecture of your computer._

#### Mailhog

Install:

```shell
go install github.com/mailhog/MailHog@latest
```

Start Mailhog:

```shell
MailHog
```

## Install Command Line Tools

### macOS 15 Sequoia and Later

```shell
npm i -g gltf-pipeline
```

```shell
brew install --cask blender
```

```shell
brew install --cask libreoffice
```

```shell
brew install \
    ocrmypdf \
    qpdf \
    exiftool \
    poppler \
    imagemagick \
    ffmpeg
```

### Debian 12 Bookworm and Later

Run [Voltaserve Conversion](conversion/README.md) with the environment variable `ENABLE_INSTALLER` set to `true`.
This will install the dependencies in the background. Incoming requests will be queued and be waiting until the installation is complete.

## Run Microservices

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Mosaic](mosaic/README.md)
- [Voltaserve Language](mosaic/README.md)
- [Voltaserve Console](console/README.md)
