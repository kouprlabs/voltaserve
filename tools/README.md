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

### Convert an DOCX file to PDF using LibreOffice

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
    "${output.*}",
    "${input}"
  ],
  "stdout": true
}
```
