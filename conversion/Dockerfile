FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-conversion

FROM mcr.microsoft.com/playwright:jammy AS runner

RUN apt-get update
RUN apt-get install software-properties-common
RUN add-apt-repository universe
RUN apt-get update

RUN apt-get install -y \
  curl \
  ffmpeg \
  gawk \
  ghostscript \
  imagemagick \
  poppler-utils \
  libimage-exiftool-perl \
  ocrmypdf \
  unzip

RUN npm i -g gltf-pipeline
RUN npm i -g @koupr/screenshot-glb

RUN apt-get install -y \
  libreoffice-base-nogui \
  libreoffice-core-nogui \
  libreoffice-writer-nogui \
  libreoffice-calc-nogui \
  libreoffice-impress-nogui \
  libreoffice-draw-nogui \
  libreoffice-math-nogui

RUN apt-get install -y \
  fonts-dejavu-core \
  fonts-dejavu-extra \
  fonts-droid-fallback \
  fonts-font-awesome \
  fonts-liberation \
  fonts-liberation2 \
  fonts-noto-cjk \
  fonts-noto-color-emoji \
  fonts-noto-core \
  fonts-noto-mono \
  fonts-ubuntu \
  fonts-ubuntu-console \
  fonts-inconsolata \
  fonts-roboto \
  fonts-roboto-slab \
  fonts-open-sans \
  fonts-source-code-pro \
  fonts-noto-sans \
  fonts-noto-serif \
  fonts-hack \
  fonts-cascadia-code \
  fonts-firacode \
  fonts-anonymous-pro \
  fonts-comfortaa \
  fonts-cantarell \
  fonts-lobster \
  fonts-meslo-nf \
  fonts-monoid \
  fonts-iosevka \
  fonts-jetbrains-mono \
  fonts-pt-sans-narrow \
  fonts-fantasque-sans-mono \
  fonts-ubuntu-monospace \
  fonts-verdana \
  fonts-gill-sans \
  fonts-futura \
  fonts-century-gothic \
  fonts-calibri \
  fonts-times-new-roman \
  fonts-arial \
  fonts-playfair-display \
  fonts-lora \
  fonts-century-schoolbook \
  fonts-bookman-old-style \
  fonts-garamond \
  fonts-georgia \
  fonts-cambria \
  fonts-palatino \
  fonts-tahoma \
  fonts-trebuchet-ms \
  fonts-baskerville

RUN apt-get install -y \
  tesseract-ocr \
  tesseract-ocr-osd \
  tesseract-ocr-ara \
  tesseract-ocr-chi-sim \
  tesseract-ocr-chi-tra \
  tesseract-ocr-eng \
  tesseract-ocr-fra \
  tesseract-ocr-deu \
  tesseract-ocr-por \
  tesseract-ocr-spa

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 8083
