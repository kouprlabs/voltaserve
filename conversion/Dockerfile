FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-conversion

FROM zenika/alpine-chrome:124-with-playwright

USER root

RUN apk update

RUN apk add --no-cache \
  curl \
  ffmpeg \
  gawk \
  ghostscript \
  imagemagick \
  poppler-utils \
  exiftool \
  perl-image-exiftool \
  ocrmypdf \
  unzip

RUN npm i -g gltf-pipeline
RUN npm i -g @shopify/screenshot-glb

RUN apk add --no-cache \
  libreoffice-base \
  libreoffice-writer \
  libreoffice-calc \
  libreoffice-impress \
  libreoffice-draw \
  libreoffice-math

RUN apk add --no-cache \
  font-dejavu \
  font-droid \
  font-awesome \
  font-liberation \
  font-liberation-sans-narrow \
  font-noto-all \
  font-ubuntu \
  font-inconsolata \
  font-roboto \
  font-opensans \
  font-hack \
  font-cantarell \
  font-jetbrains-mono \
  font-adobe-source-code-pro \
  font-adobe-utopia-100dpi \
  font-arimo \
  font-barlow \
  font-carlito \
  font-freefont \
  font-inter \
  font-montserrat \
  font-noto

RUN apk add --no-cache \
  tesseract-ocr \
  tesseract-ocr-data-osd \
  tesseract-ocr-data-ara \
  tesseract-ocr-data-chi_sim \
  tesseract-ocr-data-chi_tra \
  tesseract-ocr-data-eng \
  tesseract-ocr-data-fra \
  tesseract-ocr-data-deu \
  tesseract-ocr-data-por \
  tesseract-ocr-data-spa

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 8083