# Voltaserve Development

## Getting Started

Run infrastructure services:

```shell
docker compose up -d \
    postgres \
    minio \
    meilisearch \
    redis \
    mailhog \
    ffmpeg \
    graphicsmagick 
    libreoffice \
    ocrmypdf \
    poppler \
    tesseract
```

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Language](language/README.md)
