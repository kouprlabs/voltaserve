# Voltaserve Development

## Getting Started

### Docker

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

### SLE 15 or openSUSE Leap 15

Install dependencies:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/sle15/install.sh?t=$(date +%s)" | sh -s
```

Start infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/sle15/start.sh?t=$(date +%s)" | sh -s
```

Stop infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/sle15/start.sh?t=$(date +%s)" | sh -s
```

### RHEL 9

Install dependencies:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/rhel9/install.sh?t=$(date +%s)" | sh -s
```

Start infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/rhel9/start.sh?t=$(date +%s)" | sh -s
```

Stop infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/rhel9/stop.sh?t=$(date +%s)" | sh -s
```

### Microservices

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Language](language/README.md)
- [Voltaserve Tools](tools/README.md)
