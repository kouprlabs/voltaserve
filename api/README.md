# Voltaserve API

## Getting Started

Install [Golang](https://go.dev/doc/install).

### Install Dependencies For macOS

Install [Homebrew](https://brew.sh).

```sh
brew install ocrmypdf
brew install graphicsmagick
brew install libreoffice
brew install poppler
```

### Install Dependencies for Debian and Ubuntu

Run the follwing script:

```sh
./install-deps.sh
```

### Install Dependencies for Fedora

```sh
sudo dnf install ocrmypdf GraphicsMagick poppler-utils libreoffice ffmpeg tesseract
```

### Build and Run

Run for development:

```sh
go run .
```

Build binary:

```sh
go build .
```

Build Docker image:

```sh
docker build -t voltaserve/api .
```

## Development

For further details about development, please check this [document](./DEVELOPMENT.md).

## Documentation

[API Reference](https://voltaserve.com/api-docs/)
