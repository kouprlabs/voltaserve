# Voltaserve Tools

## Getting Started

We assume the development environment is setup as described [here](../DEVELOPMENT.md).

### Build and Run

Run for development:

```shell
air
```

Build binary:

```shell
go build .
```

Build Docker image:

```shell
docker build -t voltaserve/conversion .
```

## Example Requests

### Get Image Size with GraphicsMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "gm",
  "args": ["identify", "-format", "%w,%h", "${input}"],
  "output": true
}
```

### Convert JPEG to PNG with GraphicsMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "gm",
  "args": ["convert", "${input}", "${output.png}"],
  "stdout": true
}
```

### Generate a Thumbnail for a PDF using GraphicsMagick

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

`json`:

```json
{
  "bin": "gm",
  "args": ["convert", "-thumbnail", "x250", "${input}", "${output.png}"],
  "stdout": true
}
```

### Convert DOCX to PDF using LibreOffice

`POST http://localhost:6001/v1/run?api_key=MY_API_KEY`

**form-data:**

`file`: `image.jpg`

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

### Get TSV Data From an Image Using Tesseract

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

### Generate PDF with OCR Text Layer From an Image Using OCRmyPDF

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

### Build Docker Images

```shell
docker build -t voltaserve/ffmpeg -f Dockerfile.ffmpeg .
```

```shell
docker build -t voltaserve/graphicsmagick -f Dockerfile.graphicsmagick .
```

```shell
docker build -t voltaserve/libreoffice -f Dockerfile.libreoffice .
```

```shell
docker build -t voltaserve/ocrmypdf -f Dockerfile.ocrmypdf .
```

```shell
docker build -t voltaserve/poppler-tools -f Dockerfile.poppler-tools .
```

```shell
docker build -t voltaserve/tesseract -f Dockerfile.tesseract .
```
