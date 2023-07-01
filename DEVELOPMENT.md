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

### Using SLE 15 or openSUSE Leap 15

Run:

```shell
./scripts/sle15_install.sh
```

### Using RHEL 9

Run:

```shell
./scripts/rhel9_install.sh
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
