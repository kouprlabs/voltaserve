# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

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
  ocrmypdf \
  unzip

RUN npm i -g gltf-pipeline
RUN npm i -g @shopify/screenshot-glb

RUN apk add --no-cache \
  libreoffice-writer \
  libreoffice-calc \
  libreoffice-impress \
  libreoffice-draw \
  libreoffice-math

RUN apk add --no-cache \
  font-adobe-100dpi \
  font-adobe-100dpi-doc \
  font-adobe-75dpi \
  font-adobe-75dpi-doc \
  font-adobe-source-code-pro \
  font-adobe-utopia-100dpi \
  font-adobe-utopia-100dpi-doc \
  font-adobe-utopia-75dpi \
  font-adobe-utopia-75dpi-doc \
  font-adobe-utopia-type1 \
  font-adobe-utopia-type1-doc \
  font-alias \
  font-alias-doc \
  font-anonymous-pro-nerd \
  font-arabic-misc \
  font-arimo \
  font-arimo-nerd \
  font-awesome \
  font-awesome-brands \
  font-awesome-free \
  font-b612 \
  font-b612-mono \
  font-bakoma \
  font-bakoma-doc \
  font-bakoma-otf \
  font-bakoma-ttf \
  font-barlow \
  font-bh-100dpi \
  font-bh-100dpi-doc \
  font-bh-75dpi \
  font-bh-75dpi-doc \
  font-bitstream-100dpi \
  font-bitstream-100dpi-doc \
  font-bitstream-75dpi \
  font-bitstream-75dpi-doc \
  font-bitstream-type1 \
  font-bitstream-type1-doc \
  font-bitstrom-wera-sans-mono-nerd \
  font-cantarell \
  font-carlito \
  font-cascadia-code-nerd \
  font-comic-shanns-mono-nerd \
  font-cronyx-cyrillic \
  font-cronyx-cyrillic-doc \
  font-croscore \
  font-cursor-misc \
  font-dec-misc \
  font-degheest \
  font-dejavu \
  font-dejavu-sans-mono-nerd \
  font-droid \
  font-droid-nonlatin \
  font-droid-sans-mono-nerd \
  font-dseg \
  font-eb-garamond \
  font-fira-code-nerd \
  font-fira-mono-nerd \
  font-freefont \
  font-freefont-doc \
  font-go-mono-nerd \
  font-hack \
  font-hack-nerd \
  font-happy-times \
  font-hasklig-nerd \
  font-hermit-nerd \
  font-ia-writer-nerd \
  font-ibm-plex-mono-nerd \
  font-ibm-type1 \
  font-ibm-type1-doc \
  font-inconsolata \
  font-inconsolata-nerd \
  font-inter \
  font-iosevka \
  font-iosevka-aile \
  font-iosevka-base \
  font-iosevka-curly \
  font-iosevka-curly-slab \
  font-iosevka-slab \
  font-ipa \
  font-ipaex \
  font-isas-misc \
  font-isas-misc-doc \
  font-jetbrains-mono \
  font-jetbrains-mono-nerd \
  font-jetbrains-mono-nl \
  font-jetbrains-mono-vf \
  font-jis-misc \
  font-jis-misc-doc \
  font-karrik \
  font-liberation \
  font-liberation-mono-nerd \
  font-liberation-sans-narrow \
  font-linux-libertine \
  font-manager \
  font-manager-common \
  font-manager-doc \
  font-manager-lang \
  font-manager-nemo \
  font-manager-thunar \
  font-meslo-nerd \
  font-micro-misc \
  font-misc-cyrillic \
  font-misc-cyrillic-doc \
  font-misc-ethiopic \
  font-misc-misc \
  font-monofur-nerd \
  font-mononoki \
  font-mononoki-nerd \
  font-montserrat \
  font-mutt-misc \
  font-noto \
  font-noto-adlam \
  font-noto-ahom \
  font-noto-all \
  font-noto-arabic \
  font-noto-armenian \
  font-noto-balinese \
  font-noto-bamum \
  font-noto-bassa-vah \
  font-noto-batak \
  font-noto-bengali \
  font-noto-buginese \
  font-noto-buhid \
  font-noto-canadian-aboriginal \
  font-noto-chakma \
  font-noto-cham \
  font-noto-cherokee \
  font-noto-chorasmian \
  font-noto-cjk \
  font-noto-cjk-extra \
  font-noto-common \
  font-noto-coptic \
  font-noto-cypro-minoan \
  font-noto-devanagari \
  font-noto-dives-akuru \
  font-noto-duployan \
  font-noto-elbasan \
  font-noto-emoji \
  font-noto-ethiopic \
  font-noto-extra \
  font-noto-fangsong \
  font-noto-georgian \
  font-noto-grantha \
  font-noto-gujarati \
  font-noto-gunjala-gondi \
  font-noto-gurmukhi \
  font-noto-hanifi-rohingya \
  font-noto-hanunoo \
  font-noto-hebrew \
  font-noto-historical \
  font-noto-indic-siyaq-numbers \
  font-noto-javanese \
  font-noto-kaithi \
  font-noto-kannada \
  font-noto-kawi \
  font-noto-kayah-li \
  font-noto-khitan-small-script \
  font-noto-khmer \
  font-noto-khojki \
  font-noto-lao \
  font-noto-lepcha \
  font-noto-limbu \
  font-noto-lisu \
  font-noto-makasar \
  font-noto-malayalam \
  font-noto-masaram-gondi \
  font-noto-math \
  font-noto-mayan-numerals \
  font-noto-medefaidrin \
  font-noto-meetei-mayek \
  font-noto-mende-kikakui \
  font-noto-miao \
  font-noto-modi \
  font-noto-mongolian \
  font-noto-mro \
  font-noto-music \
  font-noto-myanmar \
  font-noto-nag-mundari \
  font-noto-nandinagari \
  font-noto-naskh-arabic \
  font-noto-nastaliq-urdu \
  font-noto-new-tai-lue \
  font-noto-newa \
  font-noto-nko \
  font-noto-nushu \
  font-noto-nyiakeng-puachue-hmong \
  font-noto-ol-chiki \
  font-noto-old-uyghur \
  font-noto-oriya \
  font-noto-osage \
  font-noto-ottoman-siyaq \
  font-noto-pahawh-hmong \
  font-noto-pau-cin-hau \
  font-noto-rashi-hebrew \
  font-noto-rejang \
  font-noto-samaritan \
  font-noto-saurashtra \
  font-noto-sharada \
  font-noto-signwriting \
  font-noto-sinhala \
  font-noto-sora-sompeng \
  font-noto-soyombo \
  font-noto-sundanese \
  font-noto-syloti-nagri \
  font-noto-symbols \
  font-noto-syriac \
  font-noto-tagbanwa \
  font-noto-tai \
  font-noto-tamil \
  font-noto-tangsa \
  font-noto-telugu \
  font-noto-test \
  font-noto-thaana \
  font-noto-thai \
  font-noto-tibetan \
  font-noto-tifinagh \
  font-noto-tirhuta \
  font-noto-toto \
  font-noto-vai \
  font-noto-vithkuqi \
  font-noto-wancho \
  font-noto-warang-citi \
  font-noto-yezidi \
  font-noto-yi \
  font-nunito \
  font-opensans \
  font-overpass \
  font-overpass-nerd \
  font-parisienne \
  font-roboto \
  font-roboto-flex \
  font-roboto-mono \
  font-schumacher-misc \
  font-screen-cyrillic \
  font-screen-cyrillic-doc \
  font-share-tech-mono-nerd \
  font-sligoil \
  font-sony-misc \
  font-source-code-pro-nerd \
  font-space-mono-nerd \
  font-sun-misc \
  font-terminus \
  font-terminus-doc \
  font-terminus-nerd \
  font-tinos-nerd \
  font-tlwg \
  font-ubuntu \
  font-ubuntu-mono-nerd \
  font-ubuntu-nerd \
  font-unifont \
  font-urw-base35 \
  font-util \
  font-util-dev \
  font-util-doc \
  font-uw-ttyp0 \
  font-victor-mono-nerd \
  font-viewer \
  font-vollkorn \
  font-winitzki-cyrillic \
  font-wqy-zenhei \
  font-xfree86-type1

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