#!/bin/bash

sudo dnf install -y git golang poppler-utils libreoffice python3-pip ghostscript tesseract wget

cd /usr/share/tesseract/tessdata
sudo wget https://github.com/kouprlabs/tessdata/raw/main/osd.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/eng.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/deu.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/fra.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/nld.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/ita.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/spa.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/por.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/swe.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/fin.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/jpn.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/chi_sim.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/chi_tra.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/kor.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/hin.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/rus.traineddata && \
sudo wget https://github.com/kouprlabs/tessdata/raw/main/ara.traineddata

sudo dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
sudo dnf install -y GraphicsMagick pngquant unpaper

sudo dnf install -y https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm
sudo dnf install -y https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-9.noarch.rpm
sudo dnf install -y dnf-plugins-core
sudo dnf config-manager --set-enabled crb
sudo dnf install -y ffmpeg --allowerasing

sudo dnf install -y automake make autoconf libtool clang
sudo dnf install -y zlib zlib-devel libjpeg libjpeg-devel libwebp libwebp-devel libtiff libtiff-devel libpng libpng-devel leptonica-devel

cd $HOME
git clone https://github.com/kouprlabs/jbig2enc.git
cd ./jbig2enc
git checkout tags/0.29
./autogen.sh
./configure --with-extra-libraries=/usr/local/lib/ --with-extra-includes=/usr/local/include/
make
sudo make install
cd $HOME
rm -rf ./jbig2enc

sudo pip3 install ocrmypdf