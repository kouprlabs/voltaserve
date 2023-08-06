# Voltaserve Development

## Getting Started

### Run Infrastructure Services With Docker

```shell
docker compose up -d \
    infra-postgres \
    infra-minio \
    infra-meilisearch \
    infra-redis \
    infra-mailhog \
    tools-ffmpeg \
    tools-imagemagick \
    tools-libreoffice \
    tools-poppler
```

### Run Microservices

Start each microservice separately in a new terminal as described here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Tools](tools/README.md)
