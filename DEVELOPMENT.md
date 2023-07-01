# Voltaserve Development

## Install Dependencies

### Using Docker

Run infrastructure services:

```shell
docker compose up -d \
    postgres \
    minio \
    meilisearch \
    redis \
    mailhog \
    exiftool \
    ffmpeg \
    imagemagick \
    libreoffice \
    ocrmypdf \
    poppler \
    tesseract
```

### Using openSUSE

Run:

```shell
./scripts/opensuse.sh
```

### Using macOS

Install [Homebrew](https://brew.sh).

Run:

```shell
./scripts/macos.sh
```

## Microservices

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Language](language/README.md)
- [Voltaserve Tools](tools/README.md)
