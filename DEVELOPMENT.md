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

### SLE / openSUSE Leap

Supported operating systems:

- [SUSE Linux Enterprise 15](https://www.suse.com/products/server)
- [openSUSE Leap 15](https://get.opensuse.org/leap)

Install dependencies:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/sle15/install.sh?t=$(date +%s)" | sh -s
```

Start infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/sle15/start.sh?t=$(date +%s)" | sh -s
```

Stop infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/sle15/start.sh?t=$(date +%s)" | sh -s
```

### RHEL

Supported operating systems:

- [Red Hat Enterprise Linux 9](https://www.redhat.com/en/technologies/linux-platforms/enterprise-linux)
- [Rocky Linux 9](https://rockylinux.org)
- [AlmaLinux 9](https://almalinux.org)
- [Oracle Linux 9](https://www.oracle.com/linux)

Install dependencies:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/rhel9/install.sh?t=$(date +%s)" | sh -s
```

Start infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/rhel9/start.sh?t=$(date +%s)" | sh -s
```

Stop infrastructure services:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/$(git symbolic-ref --short HEAD 2>/dev/null || echo 'main')/scripts/rhel9/stop.sh?t=$(date +%s)" | sh -s
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
