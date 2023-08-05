# Voltaserve Conversion

## Getting Started

Install [Golang](https://go.dev/doc/install) and [golangci-lint](https://golangci-lint.run/usage/install).

### When Running Infrastructure Services Without Docker (In SLE, openSUSE Leap or a RHEL Compatible OS)

Duplicate the file `./conversion/.env`, and rename it to `./conversion/.env.local`.

Update the following entries in the file `.env.local`:

```properties
VOLTASERVE_EXIFTOOL_URL="http://127.0.0.1:6001"
VOLTASERVE_FFMPEG_URL="http://127.0.0.1:6001"
VOLTASERVE_IMAGEMAGICK_URL="http://127.0.0.1:6001"
VOLTASERVE_LIBREOFFICE_URL="http://127.0.0.1:6001"
VOLTASERVE_OCRMYPDF_URL="http://127.0.0.1:6001"
VOLTASERVE_POPPLER_URL="http://127.0.0.1:6001"
VOLTASERVE_TESSERACT_URL="http://127.0.0.1:6001"
```

### Build and Run

Run for development:

```shell
go run .
```

Build binary:

```shell
go build .
```

### Code Linter

```shell
golangci-lint run
```

### Build Docker Image

```shell
docker build -t voltaserve/conversion .
```
