# Voltaserve Development

## Getting Started

Run infrastructure services using Docker:

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

Install dependencies on SLE 15 or openSUSE Leap 15:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/sle15/install.sh?t=$(date +%s)" | sh -s
```

Install dependencies on RHEL 9:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/scripts/rhel9/install.sh?t=$(date +%s)" | sh -s
```

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Language](language/README.md)
- [Voltaserve Tools](tools/README.md)
