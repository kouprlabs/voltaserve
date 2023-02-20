FROM golang:1.19 as builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-api

FROM ubuntu:jammy AS runner

WORKDIR /app

COPY --from=builder /build/voltaserve-api ./voltaserve-api
COPY --from=builder /build/config.yml ./config.yml
COPY --from=builder /build/templates ./templates

ENV LC_ALL=C.UTF-8
ENV LANG=C.UTF-8
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN apt-get install -y ocrmypdf
RUN apt-get install -y graphicsmagick
RUN apt-get install -y poppler-utils
RUN apt-get install -y --no-install-recommends libreoffice
RUN apt-get install -y tesseract-ocr-eng
RUN apt-get install -y tesseract-ocr-deu
RUN apt-get install -y tesseract-ocr-fra
RUN apt-get install -y tesseract-ocr-nld
RUN apt-get install -y tesseract-ocr-ita
RUN apt-get install -y tesseract-ocr-spa
RUN apt-get install -y tesseract-ocr-por
RUN apt-get install -y tesseract-ocr-jpn
RUN apt-get install -y tesseract-ocr-nor
RUN apt-get install -y tesseract-ocr-swe
RUN apt-get install -y tesseract-ocr-fin
RUN apt-get install -y tesseract-ocr-dan
RUN apt-get install -y tesseract-ocr-rus
RUN apt-get install -y tesseract-ocr-chi-sim
RUN apt-get install -y tesseract-ocr-chi-tra

ENTRYPOINT ["./voltaserve-api"]

EXPOSE 5000