# Voltaserve Conversion

## Getting Started

Install [Golang](https://go.dev/doc/install).

### Install Dependencies For macOS

Install [Homebrew](https://brew.sh).

```shell
brew install ocrmypdf
brew install graphicsmagick
brew install libreoffice
brew install poppler
brew install ffmpeg
```

### Install Dependencies for Debian and Ubuntu

Run the follwing script:

```shell
./install-deps.sh
```

### Install Dependencies for Fedora

```shell
sudo dnf install ocrmypdf GraphicsMagick poppler-utils libreoffice tesseract
```

Enable RPM Fusion repositories: <https://docs.fedoraproject.org/en-US/quick-docs/setup_rpmfusion>

```shell
sudo dnf install https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm
```

```shell
sudo dnf install https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm
```

Install `ffmpeg`:

```shell
sudo dnf swap ffmpeg-free ffmpeg --allowerasing
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

Build Docker image:

```shell
docker build -t voltaserve/conversion .
```

## Development

For further details about development, please check this [document](./DEVELOPMENT.md).

## Documentation

[API Reference](https://voltaserve.com/api-docs/)
