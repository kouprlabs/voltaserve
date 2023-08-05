# Voltaserve Tools

## Getting Started

Install [Golang](https://go.dev/doc/install) and [golangci-lint](https://golangci-lint.run/usage/install).

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

### Build Docker Images

```shell
docker build -t voltaserve/exiftoo -f ./docker/sle/Dockerfile.exiftool .
```

```shell
docker build -t voltaserve/ffmpeg -f ./docker/sle/Dockerfile.ffmpeg .
```

```shell
docker build -t voltaserve/imagemagick -f ./docker/sle/Dockerfile.imagemagick .
```

```shell
docker build -t voltaserve/libreoffice -f ./docker/sle/Dockerfile.libreoffice .
```

```shell
docker build -t voltaserve/ocrmypdf -f ./docker/sle/Dockerfile.ocrmypdf .
```

```shell
docker build -t voltaserve/poppler -f ./docker/sle/Dockerfile.poppler .
```

```shell
docker build -t voltaserve/tesseract -f D./docker/sle/ockerfile.tesseract .
```

### Example Requests

#### Get Image Size using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "identify",
  "args": ["-format", "%w,%h", "${input}"],
  "output": true
}
```

#### Convert JPEG to PNG using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "convert",
  "args": ["${input}", "${output.png}"],
  "stdout": true
}
```

#### Resize an Image using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "convert",
  "args": ["-resize", "300x", "${input}", "${output.png}"],
  "stdout": true
}
```

#### Generate a Thumbnail for a PDF using ImageMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.pdf`

`json`:

```json
{
  "bin": "convert",
  "args": ["-thumbnail", "x250", "${input}[0]", "${output.png}"],
  "stdout": true
}
```

#### Convert DOCX to PDF using LibreOffice

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.docx`

`json`:

```json
{
  "bin": "soffice",
  "args": [
    "--headless",
    "--convert-to",
    "pdf",
    "--outdir",
    "${output.*.pdf}",
    "${input}"
  ],
  "stdout": true
}
```

#### Convert PDF to Text using Poppler

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `document.pdf`

`json`:

```json
{
  "bin": "pdftotext",
  "args": ["${input}", "${output.txt}"],
  "stdout": true
}
```

#### Get TSV Data From an Image Using Tesseract

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "tesseract",
  "args": ["${input}", "${output.#.tsv}", "-l", "deu", "tsv"],
  "stdout": true
}
```

#### Generate PDF with OCR Text Layer From an Image Using OCRmyPDF

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "ocrmypdf",
  "args": [
    "--rotate-pages",
    "--clean",
    "--deskew",
    "--language=kor",
    "--image-dpi=300",
    "${input}",
    "${output}"
  ],
  "stdout": true
}
```
