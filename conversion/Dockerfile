FROM almalinux:9 as builder

WORKDIR /build

COPY . .

RUN dnf install -y golang

RUN go mod download
RUN go build -o voltaserve-conversion

FROM almalinux:9 AS runner

RUN dnf install -y git
RUN dnf install -y poppler-utils
RUN dnf install -y libreoffice
RUN dnf install -y python3-pip
RUN dnf install -y ghostscript
RUN dnf install -y tesseract
RUN dnf install -y

RUN mkdir -p /tess/traineddata

ENV TESSDATA_PREFIX=/tess/traineddata

WORKDIR /tess/traineddata

RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/eng.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/deu.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/fra.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/nld.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/ita.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/spa.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/por.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/jpn.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/nor.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/swe.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/fin.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/dan.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/rus.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/chi_sim.traineddata
RUN wget https://github.com/tesseract-ocr/tessdata/raw/main/chi_tra.traineddata

WORKDIR /build

RUN dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
RUN dnf install -y GraphicsMagick
RUN dnf install -y pngquant
RUN dnf install -y unpaper

RUN dnf install -y https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm
RUN dnf install -y https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-9.noarch.rpm
RUN dnf install -y dnf-plugins-core
RUN dnf config-manager --set-enabled crb
RUN dnf install -y ffmpeg --allowerasing

RUN dnf install -y automake 
RUN dnf install -y make
RUN dnf install -y autoconf
RUN dnf install -y libtool
RUN dnf install -y clang

RUN dnf install -y zlib
RUN dnf install -y zlib-devel
RUN dnf install -y libjpeg
RUN dnf install -y libjpeg-devel
RUN dnf install -y libwebp
RUN dnf install -y libwebp-devel
RUN dnf install -y libtiff
RUN dnf install -y libtiff-devel
RUN dnf install -y libpng
RUN dnf install -y libpng-devel
RUN dnf install -y leptonica-devel

WORKDIR /root
RUN git clone https://github.com/agl/jbig2enc.git
WORKDIR /root/jbig2enc
RUN git checkout tags/0.29
RUN ./autogen.sh
RUN ./configure --with-extra-libraries=/usr/local/lib/ --with-extra-includes=/usr/local/include/
RUN make
RUN make install
WORKDIR /root
RUN rm -rf jbig2enc

RUN pip3 install ocrmypdf

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 5001