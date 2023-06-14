#!/bin/bash

sudo dnf install -y git golang poppler-utils libreoffice python3-pip python3-devel ghostscript tesseract wget

cd /usr/share/tesseract/tessdata
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/osd.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/eng.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/deu.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/fra.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/nld.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/ita.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/spa.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/por.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/swe.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/fin.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/jpn.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/chi_sim.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/chi_tra.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/kor.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/hin.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/rus.traineddata && \
sudo wget -c https://github.com/kouprlabs/tessdata/raw/main/ara.traineddata

sudo dnf install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
sudo dnf install -y GraphicsMagick pngquant unpaper perl-Image-ExifTool

sudo dnf install -y https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm
sudo dnf install -y https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-9.noarch.rpm
sudo dnf install -y dnf-plugins-core

# For RHEL
sudo dnf config-manager --set-enabled codeready-builder-for-rhel-9-${arch}-rpms
# For Rocky Linux and AlmaLinux
sudo dnf config-manager --set-enabled crb
# For Oracle Linux
sudo dnf config-manager --set-enabled ol9_codeready_builder

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

sudo dnf module -y enable nodejs:18
sudo dnf module -y install nodejs:18/common
sudo npm install -g corepack
corepack enable