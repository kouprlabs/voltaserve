FROM golang:1.22 AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o voltaserve-conversion

FROM mcr.microsoft.com/playwright:jammy AS runner

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
  fonts-arphic-ukai \
  fonts-arphic-uming \
  fonts-beng \
  fonts-beng-extra \
  fonts-dejavu-core \
  fonts-dejavu-extra \
  fonts-deva \
  fonts-deva-extra \
  fonts-droid-fallback \
  fonts-farsiweb \
  fonts-font-awesome \
  fonts-freefont-ttf \
  fonts-gargi \
  fonts-gubbi \
  fonts-gujr \
  fonts-gujr-extra \
  fonts-guru \
  fonts-guru-extra \
  fonts-hosny-amiri \
  fonts-indic \
  fonts-kacst \
  fonts-kacst-one \
  fonts-kalapi \
  fonts-khmeros \
  fonts-khmeros-core \
  fonts-knda \
  fonts-lao \
  fonts-lato \
  fonts-liberation \
  fonts-liberation2 \
  fonts-lklug-sinhala

RUN apt-get install -y \
  fonts-lohit-beng-assamese \
  fonts-lohit-beng-bengali \
  fonts-lohit-deva \
  fonts-lohit-gujr \
  fonts-lohit-guru \
  fonts-lohit-knda \
  fonts-lohit-mlym \
  fonts-lohit-orya \
  fonts-lohit-taml \
  fonts-lohit-taml-classical \
  fonts-lohit-telu \
  fonts-manchufont \
  fonts-mathjax \
  fonts-mlym \
  fonts-nafees \
  fonts-nakula \
  fonts-navilu \
  fonts-noto-cjk \
  fonts-noto-cjk-extra \
  fonts-noto-color-emoji \
  fonts-noto-core \
  fonts-noto-mono \
  fonts-noto-ui-core \
  fonts-opensymbol \
  fonts-orya \
  fonts-orya-extra \
  fonts-pagul \
  fonts-sahadeva \
  fonts-samyak-deva \
  fonts-samyak-gujr \
  fonts-samyak-mlym \
  fonts-samyak-taml \
  fonts-sarai \
  fonts-sil-abyssinica \
  fonts-sil-ezra \
  fonts-sil-nuosusil \
  fonts-sil-padauk \
  fonts-sil-scheherazade \
  fonts-smc

RUN apt-get install -y \
  fonts-smc-anjalioldlipi \
  fonts-smc-chilanka \
  fonts-smc-dyuthi \
  fonts-smc-gayathri \
  fonts-smc-karumbi \
  fonts-smc-keraleeyam \
  fonts-smc-manjari \
  fonts-smc-meera \
  fonts-smc-rachana \
  fonts-smc-raghumalayalamsans \
  fonts-smc-suruma \
  fonts-smc-uroob \
  fonts-taml \
  fonts-telu \
  fonts-telu-extra \
  fonts-teluguvijayam \
  fonts-thai-tlwg \
  fonts-tibetan-machine \
  fonts-tlwg-garuda \
  fonts-tlwg-garuda-ttf \
  fonts-tlwg-kinnari \
  fonts-tlwg-kinnari-ttf \
  fonts-tlwg-laksaman \
  fonts-tlwg-laksaman-ttf \
  fonts-tlwg-loma \
  fonts-tlwg-loma-ttf

RUN apt-get install -y \
  fonts-tlwg-mono \
  fonts-tlwg-mono-ttf \
  fonts-tlwg-norasi \
  fonts-tlwg-norasi-ttf \
  fonts-tlwg-purisa \
  fonts-tlwg-purisa-ttf \
  fonts-tlwg-sawasdee \
  fonts-tlwg-sawasdee-ttf \
  fonts-tlwg-typewriter \
  fonts-tlwg-typewriter-ttf \
  fonts-tlwg-typist \
  fonts-tlwg-typist-ttf \
  fonts-tlwg-typo \
  fonts-tlwg-typo-ttf \
  fonts-tlwg-umpush \
  fonts-tlwg-umpush-ttf \
  fonts-tlwg-waree \
  fonts-tlwg-waree-ttf \
  fonts-ubuntu \
  fonts-ubuntu-console \
  fonts-ukij-uyghur \
  fonts-urw-base35 \
  fonts-yrsa-rasa \
  fonts-3270 \
  fonts-adf-accanthis \
  fonts-adf-baskervald \
  fonts-adf-berenis \
  fonts-adf-gillius \
  fonts-adf-ikarius \
  fonts-adf-irianis \
  fonts-adf-libris \
  fonts-adf-mekanus \
  fonts-adf-oldania \
  fonts-adf-romande

RUN apt-get install -y \
  fonts-adf-solothurn \
  fonts-adf-switzera \
  fonts-adf-tribun \
  fonts-adf-universalis \
  fonts-adf-verana \
  fonts-aenigma \
  fonts-agave \
  fonts-aksharyogini2 \
  fonts-alee \
  fonts-alegreya-sans \
  fonts-allerta \
  fonts-amiga \
  fonts-ancient-scripts \
  fonts-anonymous-pro \
  fonts-aoyagi-kouzan-t \
  fonts-aoyagi-soseki \
  fonts-apropal \
  fonts-arabeyes \
  fonts-arapey \
  fonts-arkpandora \
  fonts-arphic-bkai00mp \
  fonts-arphic-bsmi00lp \
  fonts-arphic-gbsn00lp \
  fonts-arphic-gkai00mp \
  fonts-arundina \
  fonts-atarismall \
  fonts-averia-gwf \
  fonts-averia-sans-gwf \
  fonts-averia-serif-gwf \
  fonts-b612 \
  fonts-babelstone-han \
  fonts-babelstone-modern \
  fonts-baekmuk \
  fonts-bajaderka \
  fonts-bebas-neue \
  fonts-beteckna \
  fonts-blankenburg \
  fonts-bpg-georgian \
  fonts-breip \
  fonts-bwht \
  fonts-cabin \
  fonts-cabinsketch \
  fonts-campania \
  fonts-cantarell \
  fonts-cardo \
  fonts-cascadia-code \
  fonts-cegui \
  fonts-century-catalogue \
  fonts-cherrybomb

RUN apt-get install -y \
  fonts-circos-symbols \
  fonts-clear-sans \
  fonts-cmu \
  fonts-comfortaa \
  fonts-comic-neue \
  fonts-compagnon \
  fonts-courier-prime \
  fonts-creep2 \
  fonts-croscore \
  fonts-crosextra-caladea \
  fonts-crosextra-carlito \
  fonts-cwtex-docs \
  fonts-cwtex-fs \
  fonts-cwtex-heib \
  fonts-cwtex-kai \
  fonts-cwtex-ming \
  fonts-cwtex-yen \
  fonts-dancingscript \
  fonts-ddc-uchen \
  fonts-dejavu \
  fonts-dejima-mincho \
  fonts-denemo \
  fonts-dkg-handwriting \
  fonts-dosis \
  fonts-dotgothic16 \
  fonts-dseg \
  fonts-dustin \
  fonts-dzongkha \
  fonts-ebgaramond \
  fonts-ebgaramond-extra \
  fonts-ecolier-court \
  fonts-ecolier-lignes-court \
  fonts-eeyek \
  fonts-elstob \
  fonts-elusive-icons \
  fonts-emojione \
  fonts-engadget \
  fonts-entypo \
  fonts-essays1743 \
  fonts-eurofurence \
  fonts-evertype-conakry \
  fonts-f500 \
  fonts-fantasma

RUN apt-get install -y \
  fonts-fantasque-sans \
  fonts-fanwood \
  fonts-femkeklaver \
  fonts-ferrite-core \
  fonts-firacode \
  fonts-fork-awesome \
  fonts-freefarsi \
  fonts-freefont-otf \
  fonts-gamaliel \
  fonts-gemunu-libre \
  fonts-georgewilliams \
  fonts-gfs-artemisia \
  fonts-gfs-baskerville \
  fonts-gfs-bodoni-classic \
  fonts-gfs-complutum \
  fonts-gfs-didot \
  fonts-gfs-didot-classic \
  fonts-gfs-gazis \
  fonts-gfs-neohellenic \
  fonts-gfs-olga \
  fonts-gfs-porson \
  fonts-gfs-solomos \
  fonts-gfs-theokritos \
  fonts-glasstty \
  fonts-glyphicons-halflings \
  fonts-gnutypewriter \
  fonts-go \
  fonts-gotico-antiqua \
  fonts-goudybookletter \
  fonts-hack \
  fonts-hack-otf \
  fonts-hack-ttf \
  fonts-hack-web \
  fonts-hanazono \
  fonts-havana \
  fonts-hermit \
  fonts-horai-umefont \
  fonts-hosny-thabit \
  fonts-humor-sans

RUN apt-get install -y \
  fonts-inconsolata \
  fonts-inter \
  fonts-inter-variable \
  fonts-ipaexfont \
  fonts-ipaexfont-gothic \
  fonts-ipaexfont-mincho \
  fonts-ipafont \
  fonts-ipafont-gothic \
  fonts-ipafont-mincho \
  fonts-ipamj-mincho \
  fonts-isabella \
  fonts-johnsmith-induni \
  fonts-joscelyn \
  fonts-jsmath \
  fonts-junction \
  fonts-junicode \
  fonts-jura \
  fonts-kanjistrokeorders \
  fonts-karla \
  fonts-karmilla \
  fonts-katex \
  fonts-kaushanscript \
  fonts-kiloji \
  fonts-klaudia-berenika \
  fonts-klee \
  fonts-komatuna \
  fonts-konatu \
  fonts-kouzan-mouhitsu \
  fonts-kristi \
  fonts-ldco \
  fonts-le-murmure \
  fonts-league-mono \
  fonts-league-spartan \
  fonts-leckerli-one \
  fonts-lemonada \
  fonts-levien-museum \
  fonts-levien-typoscript \
  fonts-lexi-gulim \
  fonts-lexi-saebom \
  fonts-lg-aboriginal \
  fonts-libfinal \
  fonts-lindenhill \
  fonts-linex \
  fonts-linuxlibertine \
  fonts-lmodern \
  fonts-lobster \
  fonts-lobstertwo \
  fonts-lohit-deva-marathi \
  fonts-lohit-deva-nepali \
  fonts-lyx \
  fonts-manrope \
  fonts-material-design-icons-iconfont \
  fonts-materialdesignicons-webfont \
  fonts-mathjax-extras

RUN apt-get install -y \
  fonts-meera-inimai \
  fonts-meera-taml \
  fonts-migmix \
  fonts-millimetre \
  fonts-misaki \
  fonts-mmcedar \
  fonts-mona \
  fonts-monapo \
  fonts-monlam \
  fonts-monofur \
  fonts-monoid \
  fonts-monoid-halfloose \
  fonts-monoid-halftight \
  fonts-monoid-loose \
  fonts-monoid-tight \
  fonts-motoya-l-cedar \
  fonts-motoya-l-maruberi \
  fonts-mph-2b-damase \
  fonts-mplus \
  fonts-myanmar \
  fonts-nanum \
  fonts-nanum-coding \
  fonts-nanum-eco \
  fonts-nanum-extra \
  fonts-national-park \
  fonts-naver-d2coding \
  fonts-noto \
  fonts-noto-extra \
  fonts-noto-hinted \
  fonts-noto-ui-extra \
  fonts-noto-unhinted \
  fonts-ocr-a \
  fonts-ocr-b \
  fonts-octicons \
  fonts-oflb-asana-math \
  fonts-oflb-euterpe \
  fonts-okolaks \
  fonts-oldstandard \
  fonts-open-sans \
  fonts-opendin \
  fonts-opendyslexic \
  fonts-oradano-mincho-gsrr \
  fonts-osifont \
  fonts-oxygen \
  fonts-paktype \
  fonts-paratype \
  fonts-pc \
  fonts-pc-extra \
  fonts-pecita \
  fonts-play

RUN apt-get install -y \
  fonts-povray \
  fonts-powerline \
  fonts-prociono \
  fonts-proggy \
  fonts-quattrocento \
  fonts-quicksand \
  fonts-radisnoir \
  fonts-rampart \
  fonts-recommended \
  fonts-reggae \
  fonts-ricty-diminished \
  fonts-rit-sundar \
  fonts-roadgeek \
  fonts-roboto \
  fonts-roboto-fontface \
  fonts-roboto-hinted \
  fonts-roboto-slab \
  fonts-roboto-unhinted \
  fonts-rocknroll \
  fonts-rufscript \
  fonts-sambhota-tsugring \
  fonts-sambhota-yigchung \
  fonts-samyak \
  fonts-samyak-orya \
  fonts-sawarabi-gothic \
  fonts-sawarabi-mincho \
  fonts-senamirmir-washra \
  fonts-seto \
  fonts-sil-akatab \
  fonts-sil-alkalami \
  fonts-sil-andika \
  fonts-sil-andika-compact \
  fonts-sil-andikanewbasic \
  fonts-sil-annapurna \
  fonts-sil-awami-nastaliq \
  fonts-sil-charis \
  fonts-sil-charis-compact \
  fonts-sil-dai-banna \
  fonts-sil-doulos \
  fonts-sil-doulos-compact \
  fonts-sil-galatia \
  fonts-sil-gentium \
  fonts-sil-gentium-basic \
  fonts-sil-gentiumplus \
  fonts-sil-gentiumplus-compact \
  fonts-sil-harmattan \
  fonts-sil-lateef \
  fonts-sil-mingzat \
  fonts-sil-mondulkiri \
  fonts-sil-mondulkiri-extra \
  fonts-sil-shimenkan \
  fonts-sil-shimenkan-gsm

RUN apt-get install -y \
  fonts-sil-shimenkan-guifan \
  fonts-sil-shimenkan-mas \
  fonts-sil-shimenkan-mgs \
  fonts-sil-shimenkan-salaowu \
  fonts-sil-shimenkan-sapushan \
  fonts-sil-shimenkan-taogu \
  fonts-sil-shimenkan-zonghe \
  fonts-sil-sophia-nubian \
  fonts-sil-tagmukay \
  fonts-sil-taiheritagepro \
  fonts-sil-zaghawa-beria \
  fonts-sipa-arundina \
  fonts-sixtyfour \
  fonts-sjfonts \
  fonts-solide-mirage \
  fonts-sora \
  fonts-spleen \
  fonts-staypuft \
  fonts-stick \
  fonts-stix \
  fonts-summersby \
  fonts-symbola \
  fonts-takao \
  fonts-takao-gothic \
  fonts-takao-mincho \
  fonts-takao-pgothic \
  fonts-taml-tamu \
  fonts-taml-tscu \
  fonts-terminus \
  fonts-terminus-otb \
  fonts-texgyre \
  fonts-thai-tlwg-otf \
  fonts-thai-tlwg-ttf \
  fonts-thai-tlwg-web \
  fonts-tiresias \
  fonts-tlwg-garuda-otf \
  fonts-tlwg-kinnari-otf \
  fonts-tlwg-laksaman-otf \
  fonts-tlwg-loma-otf

RUN apt-get install -y \
  fonts-tlwg-mono-otf \
  fonts-tlwg-norasi-otf \
  fonts-tlwg-purisa-otf \
  fonts-tlwg-sawasdee-otf \
  fonts-tlwg-typewriter-otf \
  fonts-tlwg-typist-otf \
  fonts-tlwg-typo-otf \
  fonts-tlwg-umpush-otf \
  fonts-tlwg-waree-otf \
  fonts-tomsontalks \
  fonts-train \
  fonts-triod-postnaja \
  fonts-tuffy \
  fonts-ubuntu-title \
  fonts-umeplus \
  fonts-umeplus-cl \
  fonts-unfonts-core \
  fonts-unfonts-extra \
  fonts-unifont \
  fonts-unikurdweb \
  fonts-uralic \
  fonts-vlgothic \
  fonts-vollkorn \
  fonts-wine \
  fonts-woowa-bm

RUN apt-get install -y \
  fonts-wqy-microhei \
  fonts-wqy-zenhei \
  fonts-yanone-kaffeesatz \
  fonts-yozvox-yozfont \
  fonts-yozvox-yozfont-antique \
  fonts-yozvox-yozfont-cute \
  fonts-yozvox-yozfont-edu \
  fonts-yozvox-yozfont-new-kana \
  fonts-yozvox-yozfont-standard-kana \
  fonts-yusei-magic \
  fonts-cns11643-kai \
  fonts-cns11643-sung \
  fonts-ibm-plex \
  fonts-ipafont-nonfree-jisx0208 \
  fonts-ipafont-nonfree-uigothic \
  fonts-jetbrains-mono \
  fonts-larabie-deco \
  fonts-larabie-straight \
  fonts-larabie-uncommon \
  fonts-mikachan \
  fonts-moe-standard-kai \
  fonts-moe-standard-song \
  fonts-mononoki \
  fonts-xfree86-nonfree \
  fonts-xfree86-nonfree-syriac

RUN apt-get install -y \
  tesseract-ocr \
  tesseract-ocr-osd \
  tesseract-ocr-ara \
  tesseract-ocr-chi-sim \
  tesseract-ocr-chi-tra \
  tesseract-ocr-nld \
  tesseract-ocr-eng \
  tesseract-ocr-fra \
  tesseract-ocr-deu \
  tesseract-ocr-hin \
  tesseract-ocr-ita \
  tesseract-ocr-jpn \
  tesseract-ocr-por \
  tesseract-ocr-rus \
  tesseract-ocr-spa \
  tesseract-ocr-swe

WORKDIR /app

COPY --from=builder /build/voltaserve-conversion ./voltaserve-conversion
COPY --from=builder /build/.env ./.env

ENTRYPOINT ["./voltaserve-conversion"]

EXPOSE 8083
