#!/bin/bash

BASE_DIR="/opt"
mkdir -p $BASE_DIR

printf_bold() {
  local msg="$1"
  printf "\e[1m%s\e[0m\n" "$msg"
}

printf_cyan() {
  local msg="$1"
  printf "\e[36m%s\e[0m\n" "$msg"
}

printf_grey() {
  local msg="$1"
  printf "\e[90m%s\e[0m\n" "$msg"
}

printf_magenta() {
  local msg="$1"
  printf "\e[35m%s\e[0m\n" "$msg"
}

printf_red() {
  local msg="$1"
  printf "\e[31m%s\e[0m\n" "$msg"
}

printf_underlined() {
  local msg="$1"
  printf "\e[4m%s\e[0m\n" "$msg"
}

is_package_installed() {
  local package_name="$1"
  local result
  result=$(zypper pa -i | grep " $package_name " | head -n 1 | awk -F "|" '{print $3}' | awk '{$1=$1};1')
  if [[ "$result" == "$package_name" ]]; then
    return 0
  else
    return 1
  fi
}

is_package_pattern_installed() {
  local package_pattern="$1"
  local expected_count="$2"
  local installed_count
  installed_count=$(zypper pa -i | grep -- "-fonts " | awk -F "|" '{print $3}' | awk '{$1=$1};1' | wc -l)
  if [[ "$installed_count" -eq "$expected_count" ]]; then
    return 0
  else
    return 1
  fi
}

is_service_running() {
  local service_name="$1"
  if sudo systemctl is-active --quiet "$service_name"; then
    return 0
  else
    return 1
  fi
}

is_brew_package_installed() {
  local package_name="$1"
  result=$(brew list --formula | grep -x "$package_name")
  if [[ "$result" == "$package_name" ]]; then
    return 0
  else
    return 1
  fi
}

install_package() {
  local package_name="$1"
  if ! is_package_installed "$package_name"; then
    printf_bold "📦  Installing package '${package_name}'..."
    sudo zypper install -y "$package_name"
    if ! is_package_installed "$package_name"; then
      printf_red "⛈️  Failed to install package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${package_name}'. Skipping."
  fi
}

install_brew_package() {
  local package_name="$1"
  local package_version="$2"
  if ! is_brew_package_installed "$package_name"; then
    printf_bold "📦  Installing brew package '${package_name}'..."
    brew install "${package_name}@${package_version}"
    if ! is_brew_package_installed "$package_name"; then
      printf_red "⛈️  Failed to install brew package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Brew package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found brew package '${package_name}'. Skipping."
  fi
}

install_pip_package() {
  local package_name="$1"
  local package_version="$2"
  local not_found="! pip show $package_name >/dev/null 2>&1"
  if eval "$not_found"; then
    printf_bold "🐍  Installing Python package '${package_name}'..."
    pip3 install "${package_name}==${package_version}"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install Python package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Python package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found Python package '${package_name}'. Skipping."
  fi
}

install_corepack() {
  local not_found="! npm list -g corepack >/dev/null 2>&1"
  if eval "$not_found"; then
    printf_bold "💎  Installing NPM package 'corepack'..."
    sudo npm install -g corepack@0.18.1
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install NPM package 'corepack'. Aborting."
      exit 1
    else
      printf_bold "✅  NPM package 'corepack' installed successfully."
    fi
  else
    printf_bold "✅  Found NPM package 'corepack'. Skipping."
  fi
}

install_golangci() {
  local golangci_bin="${HOME}/bin/golangci-lint"
  local not_found="! (command -v $golangci_bin >/dev/null 2>&1 && $golangci_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "🐹  Installing Go binary '${golangci_bin}'..."
    mkdir -p "$HOME/bin"
    cd "$HOME" || exit
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.2
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install Go binary '${golangci_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Go binary '${golangci_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found Go binary '${golangci_bin}'. Skipping."
  fi
}

install_swag() {
  local swag_bin="${HOME}/bin/swag"
  local not_found="! (command -v $swag_bin >/dev/null 2>&1 && $swag_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "🐹  Installing Go binary '${swag_bin}'..."
    go install github.com/swaggo/swag/cmd/swag@v1.8.12
    mkdir -p "$HOME/bin"
    mv "$(go env GOPATH)/bin/swag" "$HOME/bin/swag"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install Go binary '${swag_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Go binary '${swag_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found Go binary '${swag_bin}'. Skipping."
  fi
}

install_air() {
  local air_bin="${HOME}/bin/air"
  local not_found="! (command -v $air_bin >/dev/null 2>&1 && $air_bin -v >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "🐹  Installing Go binary '${air_bin}'..."
    go install github.com/cosmtrek/air@v1.44.0
    mkdir -p "$HOME/bin"
    mv "$(go env GOPATH)/bin/air" "$HOME/bin/air"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install Go binary '${air_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Go binary '${air_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found Go binary '${air_bin}'. Skipping."
  fi
}

install_postgres() {
  local postgres_service="postgresql"
  local not_found='! systemctl list-unit-files | grep -q '"${postgres_service}.service"''
  if eval "$not_found"; then
    printf_bold "📦  Installing service '${postgres_service}'..."
    sudo zypper install -y postgresql15-server
    sudo systemctl enable --now postgresql
    sudo su postgres <<EOF
psql -c "CREATE USER voltaserve WITH PASSWORD 'voltaserve';"
psql -c "CREATE DATABASE voltaserve;"
psql -c "GRANT ALL ON DATABASE voltaserve TO voltaserve;"
psql -c "ALTER DATABASE voltaserve OWNER TO voltaserve;"
exit
EOF
    sudo su postgres <<EOF
psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"
exit
EOF
    sudo sed -i 's/peer/md5/g' /var/lib/pgsql/data/pg_hba.conf
    sudo sed -i 's/ident/md5/g' /var/lib/pgsql/data/pg_hba.conf
    sudo systemctl restart postgresql
    git_branch=$(git symbolic-ref --short HEAD 2>/dev/null)
    curl -fsSL "https://raw.githubusercontent.com/kouprlabs/voltaserve/${git_branch}/postgres/schema.sql" | PGPASSWORD="voltaserve" psql -U "voltaserve"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install service '${postgres_service}'. Aborting."
      exit 1
    else
      printf_bold "✅  Service '${postgres_service}' installed successfully."
    fi
  else
    printf_bold "✅  Found service '$postgres_service'. Skipping."
  fi
}

install_redis() {
  local package_name="redis7"
  if ! is_package_installed "$package_name"; then
    printf_bold "📦  Installing package '${package_name}'..."
    sudo zypper install -y redis7
    echo "[Unit]
Description=Redis Server
After=network.target

[Service]
ExecStart=/usr/sbin/redis-server
Restart=always

[Install]
WantedBy=default.target" | sudo tee /etc/systemd/system/redis.service >/dev/null
    sudo systemctl daemon-reload
    sudo systemctl enable --now redis.service
    if ! is_package_installed "$package_name" || ! is_service_running "redis.service"; then
      printf_red "⛈️  Failed to install package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${package_name}'. Skipping."
  fi
}

install_minio() {
  local minio_pkg="minio"
  local not_found="! rpm -q $minio_pkg >/dev/null"
  if eval "$not_found"; then
    printf_bold "📦  Installing package '${minio_pkg}'..."
    local minio_rpm="minio-20230629051228.0.0.x86_64.rpm"
    sudo wget -c "https://dl.min.io/server/minio/release/linux-amd64/${minio_rpm}" -P $BASE_DIR
    sudo zypper --no-gpg-checks install -y "${BASE_DIR}/${minio_rpm}"
    sudo rm -f "${BASE_DIR}/${minio_rpm}"
    sudo mkdir -p "${BASE_DIR}/minio"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install package '${minio_pkg}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${minio_pkg}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${minio_pkg}'. Skipping."
  fi
}

install_meilisearch() {
  local meilisearch_bin="${BASE_DIR}/meilisearch/meilisearch"
  local not_found="! (command -v $meilisearch_bin >/dev/null 2>&1 && $meilisearch_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "📦  Installing binary '${meilisearch_bin}'..."
    sudo mkdir -p "${BASE_DIR}/meilisearch"
    cd "${BASE_DIR}/meilisearch" || exit
    sudo wget -c "https://github.com/meilisearch/meilisearch/releases/download/v1.2.0/meilisearch-linux-amd64"
    sudo mv ./meilisearch-linux-amd64 ./meilisearch
    sudo chmod +x $meilisearch_bin
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install binary '${meilisearch_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Binary '${meilisearch_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found binary '${meilisearch_bin}'. Skipping."
  fi
}

install_mailhog() {
  local mailhog_bin="${BASE_DIR}/mailhog/MailHog_linux_amd64"
  local not_found="! (command -v $mailhog_bin >/dev/null 2>&1 && $mailhog_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "📦  Installing binary '${mailhog_bin}'..."
    sudo mkdir -p "${BASE_DIR}/mailhog"
    sudo wget -c https://github.com/mailhog/MailHog/releases/download/v1.0.1/MailHog_linux_amd64 -P "${BASE_DIR}/mailhog"
    sudo chmod +x $mailhog_bin
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install binary '${mailhog_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Binary '${mailhog_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found binary '${mailhog_bin}'. Skipping."
  fi
}

install_tesseract() {
  local package_name="tesseract-ocr"
  if ! is_package_installed "$package_name"; then
    printf_bold "📦  Installing package '${package_name}'..."
    sudo zypper install -y \
      tesseract-ocr \
      tesseract-ocr-traineddata-osd \
      tesseract-ocr-traineddata-eng \
      tesseract-ocr-traineddata-deu \
      tesseract-ocr-traineddata-fra \
      tesseract-ocr-traineddata-nld \
      tesseract-ocr-traineddata-ita \
      tesseract-ocr-traineddata-spa \
      tesseract-ocr-traineddata-por \
      tesseract-ocr-traineddata-swe \
      tesseract-ocr-traineddata-fin \
      tesseract-ocr-traineddata-jpn \
      tesseract-ocr-traineddata-chi_sim \
      tesseract-ocr-traineddata-chi_tra \
      tesseract-ocr-traineddata-hin \
      tesseract-ocr-traineddata-rus \
      tesseract-ocr-traineddata-ara
    if ! is_package_installed "$package_name"; then
      printf_red "⛈️  Failed to install package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${package_name}'. Skipping."
  fi
}

install_libreoffice() {
  local package_name="libreoffice"
  if ! is_package_installed $package_name; then
    printf_bold "📦  Installing package '${package_name}'..."
    sudo zypper install -y \
      libreoffice \
      libreoffice-writer \
      libreoffice-calc \
      libreoffice-impress \
      libreoffice-draw \
      libreoffice-math \
      libreoffice-base
    if ! is_package_installed "$package_name"; then
      printf_red "⛈️  Failed to install package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${package_name}'. Skipping."
  fi
}

install_fonts() {
  local package_pattern="*-fonts"
  local package_count=863
  if ! is_package_pattern_installed "$package_pattern" $package_count; then
    printf_bold "📦  Installing fonts..."
    sudo zypper install -y \
      adinatha-fonts \
      adobe-sourcecodepro-fonts \
      adobe-sourcehansans-cn-fonts \
      adobe-sourcehansans-hk-fonts \
      adobe-sourcehansans-jp-fonts \
      adobe-sourcehansans-kr-fonts \
      adobe-sourcehansans-tw-fonts \
      adobe-sourcehanserif-cn-fonts \
      adobe-sourcehanserif-hk-fonts \
      adobe-sourcehanserif-jp-fonts \
      adobe-sourcehanserif-kr-fonts \
      adobe-sourcehanserif-tw-fonts \
      adobe-sourcesans3-fonts \
      adobe-sourcesanspro-fonts \
      adobe-sourceserif4-fonts \
      adobe-sourceserifpro-fonts \
      aldusleaf-crimson-text-fonts \
      alee-fonts \
      aliftype-amiri-fonts \
      antijingoist-opendyslexic-fonts \
      arabic-ae-fonts \
      arabic-amiri-fonts \
      arabic-bitmap-fonts \
      arabic-fonts \
      arabic-kacst-fonts \
      arabic-kacstone-fonts \
      arabic-naqsh-fonts \
      arphic-bkai00mp-fonts \
      arphic-bsmi00lp-fonts \
      arphic-fonts \
      arphic-gbsn00lp-fonts \
      arphic-gkai00mp-fonts \
      arphic-ukai-fonts \
      arphic-uming-fonts \
      autonym-fonts \
      avesta-fonts \
      babelstone-han-fonts \
      babelstone-marchen-fonts \
      babelstone-modern-fonts \
      babelstone-ogham-fonts \
      babelstone-phags-pa-fonts \
      babelstone-runic-fonts \
      baekmuk-bitmap-fonts \
      baekmuk-ttf-fonts \
      bitstream-vera-fonts \
      blockzone-fonts \
      bpg-fonts \
      cadsondemak-fonts \
      cantarell-fonts \
      caslon-fonts \
      cm-unicode-fonts \
      consoleet-darwin-fonts \
      consoleet-fixedsys-fonts \
      consoleet-kbd-fonts \
      consoleet-oldschoolpc-fonts \
      consoleet-terminus-fonts \
      consoleet-xorg-fonts \
      courier-prime-fonts \
      cozette-fonts \
      cpmono_v07-fonts \
      culmus-ancient-semitic-fonts \
      culmus-fonts \
      cyreal-alice-fonts \
      cyreal-junge-fonts \
      cyreal-lobster-cyrillic-fonts \
      cyreal-lora-fonts \
      cyreal-marko-horobchyk-fonts \
      cyreal-marmelad-fonts \
      cyreal-wire-fonts \
      dai-banna-fonts \
      dejavu-fonts \
      delaguardo-inconsolata_lgc-fonts \
      dina-bitmap-fonts \
      eb-garamond-fonts \
      eeyek-fonts \
      efont-serif-fonts \
      efont-unicode-bitmap-fonts \
      farsi-fonts \
      finalcut-bitmap-fonts \
      fira-code-fonts \
      fontawesome-fonts \
      free-ttf-fonts \
      gdouros-abydos-fonts \
      gdouros-aegean-fonts \
      gdouros-aegyptus-fonts \
      gdouros-akkadian-fonts \
      gdouros-alfios-fonts \
      gdouros-analecta-fonts \
      gdouros-anatolian-fonts \
      gdouros-atavyros-fonts \
      gdouros-maya-fonts \
      gdouros-musica-fonts \
      gdouros-symbola-fonts \
      gdouros-text-fonts \
      gdouros-unidings-fonts \
      gnu-free-fonts \
      gnu-unifont-bitmap-fonts \
      gnu-unifont-legacy-bitmap-fonts \
      google-alegreya-fonts \
      google-alegreya-sans-fonts \
      google-allerta-fonts \
      google-anonymouspro-fonts \
      google-arimo-fonts \
      google-cabin-fonts \
      google-caladea-fonts \
      google-cardo-fonts \
      google-carlito-fonts \
      google-cousine-fonts \
      google-croscore-fonts \
      google-droid-fonts \
      google-exo-fonts \
      google-inconsolata-fonts \
      google-lekton-fonts \
      google-merriweather-fonts \
      google-nobile-fonts \
      google-opensans-fonts \
      google-poppins-fonts \
      google-roboto-fonts \
      google-roboto-mono-fonts \
      google-tinos-fonts \
      google-worksans-fonts \
      hack-fonts \
      hartke-aurulentsans-fonts \
      indic-fonts \
      inter-fonts \
      intlfonts-arabic-bitmap-fonts \
      intlfonts-asian-bitmap-fonts \
      intlfonts-bdf-fonts \
      intlfonts-chinese-big-bitmap-fonts \
      intlfonts-chinese-bitmap-fonts \
      intlfonts-ethiopic-bitmap-fonts \
      intlfonts-euro-bitmap-fonts \
      intlfonts-japanese-big-bitmap-fonts \
      intlfonts-japanese-bitmap-fonts \
      intlfonts-phonetic-bitmap-fonts \
      intlfonts-ttf-fonts \
      intlfonts-type1-fonts \
      iosevka-aile-fonts \
      iosevka-curly-fonts \
      iosevka-curly-slab-fonts \
      iosevka-etoile-fonts \
      iosevka-fonts \
      iosevka-slab-fonts \
      iosevka-ss01-fonts \
      iosevka-ss02-fonts \
      iosevka-ss03-fonts \
      iosevka-ss04-fonts \
      iosevka-ss05-fonts \
      iosevka-ss06-fonts \
      iosevka-ss07-fonts \
      iosevka-ss08-fonts \
      iosevka-ss09-fonts \
      iosevka-ss10-fonts \
      iosevka-ss11-fonts \
      iosevka-ss12-fonts \
      iosevka-ss13-fonts \
      iosevka-ss14-fonts \
      iosevka-ss15-fonts \
      iosevka-ss16-fonts \
      iosevka-ss17-fonts \
      iosevka-ss18-fonts \
      ipa-ex-gothic-fonts \
      ipa-ex-mincho-fonts \
      ipa-gothic-bold-fonts \
      ipa-gothic-bolditalic-fonts \
      ipa-gothic-fonts \
      ipa-gothic-italic-fonts \
      ipa-mincho-fonts \
      ipa-pgothic-bold-fonts \
      ipa-pgothic-bolditalic-fonts \
      ipa-pgothic-fonts \
      ipa-pgothic-italic-fonts \
      ipa-pmincho-fonts \
      ipa-uigothic-fonts \
      jetbrains-mono-fonts \
      jomolhari-fonts \
      js-technology-fonts \
      kde-oxygen-fonts \
      khmeros-fonts \
      kika-fixedsys-fonts \
      lato-fonts \
      liberation-fonts \
      libertinus-fonts \
      lilypond-emmentaler-fonts \
      linux-libertine-fonts \
      lklug-fonts \
      lomt-blackout-fonts \
      lomt-chunk-fonts \
      lomt-fanwood-fonts \
      lomt-goudybookletter-fonts \
      lomt-junction-fonts \
      lomt-knewave-fonts \
      lomt-leaguegothic-fonts \
      lomt-lindenhill-fonts \
      lomt-orbitron-fonts \
      lomt-ostrichsans-fonts \
      lomt-prociono-fonts \
      lomt-script1-fonts \
      lomt-sniglet-fonts \
      lomt-sortsmillgoudy-fonts \
      lyx-fonts \
      manchu-fonts \
      mathgl-fonts \
      mathjax-ams-fonts \
      mathjax-caligraphic-fonts \
      mathjax-fraktur-fonts \
      mathjax-main-fonts \
      mathjax-math-fonts \
      mathjax-sansserif-fonts \
      mathjax-script-fonts \
      mathjax-size1-fonts \
      mathjax-size2-fonts \
      mathjax-size3-fonts \
      mathjax-size4-fonts \
      mathjax-typewriter-fonts \
      mathjax-winchrome-fonts \
      mathjax-winie6-fonts \
      meslo-lg-fonts \
      mgopen-fonts \
      miao-fonts \
      mikachan-fonts \
      mingzat-fonts \
      monapo-fonts \
      mongolian-fonts \
      motoya-lcedar-fonts \
      motoya-lmaru-fonts \
      mplus-code-latin-variable-fonts \
      mplus-code-latin50-fonts \
      mplus-code-latin60-fonts \
      mplus-fonts \
      mplus1-code-fonts \
      mplus1-code-variable-fonts \
      mplus1-fonts \
      mplus1-variable-fonts \
      mplus2-fonts \
      mplus2-variable-fonts \
      musescore-fonts \
      namdhinggo-fonts \
      nanum-fonts \
      nanum-gothic-coding-fonts \
      noto-coloremoji-fonts \
      noto-kufiarabic-fonts \
      noto-mono-fonts \
      noto-naskharabic-fonts \
      noto-naskharabic-ui-fonts \
      noto-nastaliqurdu-fonts \
      noto-sans-adlam-fonts \
      noto-sans-adlamunjoined-fonts \
      noto-sans-anatolianhieroglyphs-fonts \
      noto-sans-arabic-fonts \
      noto-sans-arabic-ui-fonts \
      noto-sans-armenian-fonts \
      noto-sans-avestan-fonts \
      noto-sans-balinese-fonts \
      noto-sans-bamum-fonts \
      noto-sans-batak-fonts \
      noto-sans-bengali-fonts \
      noto-sans-bengali-ui-fonts \
      noto-sans-brahmi-fonts \
      noto-sans-buginese-fonts \
      noto-sans-buhid-fonts \
      noto-sans-canadianaboriginal-fonts \
      noto-sans-carian-fonts \
      noto-sans-chakma-fonts \
      noto-sans-cham-fonts \
      noto-sans-cherokee-fonts \
      noto-sans-cjk-fonts \
      noto-sans-coptic-fonts \
      noto-sans-cuneiform-fonts \
      noto-sans-cypriot-fonts \
      noto-sans-deseret-fonts \
      noto-sans-devanagari-fonts \
      noto-sans-devanagari-ui-fonts \
      noto-sans-display-fonts \
      noto-sans-egyptianhieroglyphs-fonts \
      noto-sans-ethiopic-fonts \
      noto-sans-fonts \
      noto-sans-georgian-fonts \
      noto-sans-glagolitic-fonts \
      noto-sans-gothic-fonts \
      noto-sans-gujarati-fonts \
      noto-sans-gujarati-ui-fonts \
      noto-sans-gurmukhi-fonts \
      noto-sans-gurmukhi-ui-fonts \
      noto-sans-hanunoo-fonts \
      noto-sans-hebrew-fonts \
      noto-sans-imperialaramaic-fonts \
      noto-sans-inscriptionalpahlavi-fonts \
      noto-sans-inscriptionalparthian-fonts \
      noto-sans-javanese-fonts \
      noto-sans-jp-black-fonts \
      noto-sans-jp-bold-fonts \
      noto-sans-jp-demilight-fonts \
      noto-sans-jp-fonts \
      noto-sans-jp-light-fonts \
      noto-sans-jp-medium-fonts \
      noto-sans-jp-mono-fonts \
      noto-sans-jp-regular-fonts \
      noto-sans-jp-thin-fonts \
      noto-sans-kaithi-fonts \
      noto-sans-kannada-fonts \
      noto-sans-kannada-ui-fonts \
      noto-sans-kayahli-fonts \
      noto-sans-kharoshthi-fonts \
      noto-sans-khmer-fonts \
      noto-sans-khmer-ui-fonts \
      noto-sans-kr-black-fonts \
      noto-sans-kr-bold-fonts \
      noto-sans-kr-demilight-fonts \
      noto-sans-kr-fonts \
      noto-sans-kr-light-fonts \
      noto-sans-kr-medium-fonts \
      noto-sans-kr-mono-fonts \
      noto-sans-kr-regular-fonts \
      noto-sans-kr-thin-fonts \
      noto-sans-lao-fonts \
      noto-sans-lao-ui-fonts \
      noto-sans-lepcha-fonts \
      noto-sans-limbu-fonts \
      noto-sans-linearb-fonts \
      noto-sans-lisu-fonts \
      noto-sans-lycian-fonts \
      noto-sans-lydian-fonts \
      noto-sans-malayalam-fonts \
      noto-sans-malayalam-ui-fonts \
      noto-sans-mandaic-fonts \
      noto-sans-meeteimayek-fonts \
      noto-sans-mongolian-fonts \
      noto-sans-mono-fonts \
      noto-sans-myanmar-fonts \
      noto-sans-myanmar-ui-fonts \
      noto-sans-newtailue-fonts \
      noto-sans-nko-fonts \
      noto-sans-ogham-fonts \
      noto-sans-olchiki-fonts \
      noto-sans-olditalic-fonts \
      noto-sans-oldpersian-fonts \
      noto-sans-oldsoutharabian-fonts \
      noto-sans-oldturkic-fonts \
      noto-sans-oriya-fonts \
      noto-sans-oriya-ui-fonts \
      noto-sans-osage-fonts \
      noto-sans-osmanya-fonts \
      noto-sans-phagspa-fonts \
      noto-sans-phoenician-fonts \
      noto-sans-rejang-fonts \
      noto-sans-runic-fonts \
      noto-sans-samaritan-fonts \
      noto-sans-saurashtra-fonts \
      noto-sans-sc-black-fonts \
      noto-sans-sc-bold-fonts \
      noto-sans-sc-demilight-fonts \
      noto-sans-sc-fonts \
      noto-sans-sc-light-fonts \
      noto-sans-sc-medium-fonts \
      noto-sans-sc-mono-fonts \
      noto-sans-sc-regular-fonts \
      noto-sans-sc-thin-fonts \
      noto-sans-shavian-fonts \
      noto-sans-sinhala-fonts \
      noto-sans-sinhala-ui-fonts \
      noto-sans-sundanese-fonts \
      noto-sans-sylotinagri-fonts \
      noto-sans-symbols-fonts \
      noto-sans-symbols2-fonts \
      noto-sans-syriaceastern-fonts \
      noto-sans-syriacestrangela-fonts \
      noto-sans-syriacwestern-fonts \
      noto-sans-tagalog-fonts \
      noto-sans-tagbanwa-fonts \
      noto-sans-taile-fonts \
      noto-sans-taitham-fonts \
      noto-sans-taiviet-fonts \
      noto-sans-tamil-fonts \
      noto-sans-tamil-ui-fonts \
      noto-sans-tc-black-fonts \
      noto-sans-tc-bold-fonts \
      noto-sans-tc-demilight-fonts \
      noto-sans-tc-fonts \
      noto-sans-tc-light-fonts \
      noto-sans-tc-medium-fonts \
      noto-sans-tc-mono-fonts \
      noto-sans-tc-regular-fonts \
      noto-sans-tc-thin-fonts \
      noto-sans-telugu-fonts \
      noto-sans-telugu-ui-fonts \
      noto-sans-thaana-fonts \
      noto-sans-thai-fonts \
      noto-sans-thai-ui-fonts \
      noto-sans-tibetan-fonts \
      noto-sans-tifinagh-fonts \
      noto-sans-ugaritic-fonts \
      noto-sans-vai-fonts \
      noto-sans-yi-fonts \
      noto-serif-armenian-fonts \
      noto-serif-bengali-fonts \
      noto-serif-devanagari-fonts \
      noto-serif-display-fonts \
      noto-serif-ethiopic-fonts \
      noto-serif-fonts \
      noto-serif-georgian-fonts \
      noto-serif-gujarati-fonts \
      noto-serif-hebrew-fonts \
      noto-serif-jp-black-fonts \
      noto-serif-jp-bold-fonts \
      noto-serif-jp-extralight-fonts \
      noto-serif-jp-fonts \
      noto-serif-jp-light-fonts \
      noto-serif-jp-medium-fonts \
      noto-serif-jp-regular-fonts \
      noto-serif-jp-semibold-fonts \
      noto-serif-kannada-fonts \
      noto-serif-khmer-fonts \
      noto-serif-kr-black-fonts \
      noto-serif-kr-bold-fonts \
      noto-serif-kr-extralight-fonts \
      noto-serif-kr-fonts \
      noto-serif-kr-light-fonts \
      noto-serif-kr-medium-fonts \
      noto-serif-kr-regular-fonts \
      noto-serif-kr-semibold-fonts \
      noto-serif-lao-fonts \
      noto-serif-malayalam-fonts \
      noto-serif-myanmar-fonts \
      noto-serif-sc-black-fonts \
      noto-serif-sc-bold-fonts \
      noto-serif-sc-extralight-fonts \
      noto-serif-sc-fonts \
      noto-serif-sc-light-fonts \
      noto-serif-sc-medium-fonts \
      noto-serif-sc-regular-fonts \
      noto-serif-sc-semibold-fonts \
      noto-serif-sinhala-fonts \
      noto-serif-tamil-fonts \
      noto-serif-tc-black-fonts \
      noto-serif-tc-bold-fonts \
      noto-serif-tc-extralight-fonts \
      noto-serif-tc-fonts \
      noto-serif-tc-light-fonts \
      noto-serif-tc-medium-fonts \
      noto-serif-tc-regular-fonts \
      noto-serif-tc-semibold-fonts \
      noto-serif-telugu-fonts \
      noto-serif-thai-fonts \
      nuosu-fonts \
      officecodepro-fonts \
      opengost-otf-fonts \
      opengost-ttf-fonts \
      orkhon-fonts \
      paglinawan-quicksand-fonts \
      pagul-fonts \
      paratype-pt-mono-fonts \
      paratype-pt-sans-fonts \
      paratype-pt-serif-fonts \
      patterns-fonts-fonts \
      phetsarath-fonts \
      pothana2000-fonts \
      powerline-fonts \
      raleway-fonts \
      rmit-sansforgetica-fonts \
      rovasiras-kende-fonts \
      rovasiras-roga-fonts \
      saja-cascadia-code-fonts \
      saweri-fonts \
      sazanami-fonts \
      sgi-bitmap-fonts \
      shannpersand-comicshanns-fonts \
      sil-abyssinica-fonts \
      sil-andika-fonts \
      sil-charis-fonts \
      sil-doulos-fonts \
      sil-gentium-fonts \
      sil-mondulkiri-fonts \
      sil-padauk-fonts \
      steinberg-bravura-fonts \
      stix-fonts \
      stix-integrals-fonts \
      stix-pua-fonts \
      stix-sizes-fonts \
      stix-variants-fonts \
      sundanese-unicode-fonts \
      tagbanwa-fonts \
      tai-heritage-pro-fonts \
      terminus-bitmap-fonts \
      terminus-ttf-fonts \
      termsyn-bitmap-fonts \
      texlive-academicons-fonts \
      texlive-accanthis-fonts \
      texlive-adforn-fonts \
      texlive-adfsymbols-fonts \
      texlive-aesupp-fonts \
      texlive-alegreya-fonts \
      texlive-alfaslabone-fonts \
      texlive-algolrevived-fonts \
      texlive-alkalami-fonts \
      texlive-allrunes-fonts \
      texlive-almendra-fonts \
      texlive-almfixed-fonts \
      texlive-amiri-fonts \
      texlive-amsfonts-fonts \
      texlive-anonymouspro-fonts \
      texlive-antiqua-fonts \
      texlive-antt-fonts \
      texlive-arabi-fonts \
      texlive-arabtex-fonts \
      texlive-aramaic-serto-fonts \
      texlive-archaic-fonts \
      texlive-archivo-fonts \
      texlive-arev-fonts \
      texlive-arimo-fonts \
      texlive-armtex-fonts \
      texlive-arphic-fonts \
      texlive-arphic-ttf-fonts \
      texlive-arvo-fonts \
      texlive-Asana-Math-fonts \
      texlive-asapsym-fonts \
      texlive-ascii-font-fonts \
      texlive-ascmac-fonts \
      texlive-aspectratio-fonts \
      texlive-atkinson-fonts \
      texlive-augie-fonts \
      texlive-auncial-new-fonts \
      texlive-aurical-fonts \
      texlive-avantgar-fonts \
      texlive-baekmuk-fonts \
      texlive-bartel-chess-fonts \
      texlive-baskervald-fonts \
      texlive-baskervaldx-fonts \
      texlive-baskervillef-fonts \
      texlive-bbold-type1-fonts \
      texlive-belleek-fonts \
      texlive-bera-fonts \
      texlive-berenisadf-fonts \
      texlive-beuron-fonts \
      texlive-bguq-fonts \
      texlive-bitter-fonts \
      texlive-bookhands-fonts \
      texlive-bookman-fonts \
      texlive-boondox-fonts \
      texlive-brushscr-fonts \
      texlive-burmese-fonts \
      texlive-cabin-fonts \
      texlive-caladea-fonts \
      texlive-calligra-type1-fonts \
      texlive-cantarell-fonts \
      texlive-carlito-fonts \
      texlive-carolmin-ps-fonts \
      texlive-cascadia-code-fonts \
      texlive-cbcoptic-fonts \
      texlive-cbfonts-fonts \
      texlive-cc-pl-fonts \
      texlive-ccicons-fonts \
      texlive-charissil-fonts \
      texlive-charter-fonts \
      texlive-chemarrow-fonts \
      texlive-chivo-fonts \
      texlive-cinzel-fonts \
      texlive-cjhebrew-fonts \
      texlive-clara-fonts \
      texlive-clearsans-fonts \
      texlive-cm-lgc-fonts \
      texlive-cm-super-fonts \
      texlive-cm-unicode-fonts \
      texlive-cmathbb-fonts \
      texlive-cmcyr-fonts \
      texlive-cmexb-fonts \
      texlive-cmll-fonts \
      texlive-cmsrb-fonts \
      texlive-cmupint-fonts \
      texlive-cochineal-fonts \
      texlive-coelacanth-fonts \
      texlive-comfortaa-fonts \
      texlive-comicneue-fonts \
      texlive-concmath-fonts \
      texlive-context-fonts \
      texlive-cormorantgaramond-fonts \
      texlive-countriesofeurope-fonts \
      texlive-courier-fonts \
      texlive-courierten-fonts \
      texlive-crimson-fonts \
      texlive-crimsonpro-fonts \
      texlive-cryst-fonts \
      texlive-cs-fonts \
      texlive-cuprum-fonts \
      texlive-cyklop-fonts \
      texlive-dad-fonts \
      texlive-dantelogo-fonts \
      texlive-dejavu-fonts \
      texlive-dictsym-fonts \
      texlive-domitian-fonts \
      texlive-doublestroke-fonts \
      texlive-doulossil-fonts \
      texlive-dozenal-fonts \
      texlive-drm-fonts \
      texlive-droid-fonts \
      texlive-dsserif-fonts \
      texlive-dutchcal-fonts \
      texlive-ebgaramond-fonts \
      texlive-eczar-fonts \
      texlive-ektype-tanka-fonts \
      texlive-electrum-fonts \
      texlive-epigrafica-fonts \
      texlive-epiolmec-fonts \
      texlive-erewhon-fonts \
      texlive-erewhon-math-fonts \
      texlive-esint-type1-fonts \
      texlive-esrelation-fonts \
      texlive-esstix-fonts \
      texlive-esvect-fonts \
      texlive-etbb-fonts \
      texlive-ethiop-t1-fonts \
      texlive-eurosym-fonts \
      texlive-fandol-fonts \
      texlive-fbb-fonts \
      texlive-fdsymbol-fonts \
      texlive-fetamont-fonts \
      texlive-fge-fonts \
      texlive-figbas-fonts \
      texlive-fira-fonts \
      texlive-firamath-fonts \
      texlive-foekfont-fonts \
      texlive-fonetika-fonts \
      texlive-fontawesome-fonts \
      texlive-fontawesome5-fonts \
      texlive-fontmfizz-fonts \
      texlive-fonts-churchslavonic-fonts \
      texlive-fonts-tlwg-fonts \
      texlive-forum-fonts \
      texlive-fourier-fonts \
      texlive-fpl-fonts \
      texlive-frcursive-fonts \
      texlive-frederika2016-fonts \
      texlive-frimurer-fonts \
      texlive-garamond-libre-fonts \
      texlive-garamond-math-fonts \
      texlive-gentium-tug-fonts \
      texlive-gfsartemisia-fonts \
      texlive-gfsbaskerville-fonts \
      texlive-gfsbodoni-fonts \
      texlive-gfscomplutum-fonts \
      texlive-gfsdidot-fonts \
      texlive-gfsdidotclassic-fonts \
      texlive-gfsneohellenic-fonts \
      texlive-gfsneohellenicmath-fonts \
      texlive-gfsporson-fonts \
      texlive-gfssolomos-fonts \
      texlive-gillius-fonts \
      texlive-gnu-freefont-fonts \
      texlive-gofonts-fonts \
      texlive-gregoriotex-fonts \
      texlive-grotesq-fonts \
      texlive-gudea-fonts \
      texlive-hacm-fonts \
      texlive-haranoaji-extra-fonts \
      texlive-haranoaji-fonts \
      texlive-helmholtz-ellis-ji-notation-fonts \
      texlive-helvetic-fonts \
      texlive-heuristica-fonts \
      texlive-hfbright-fonts \
      texlive-hindmadurai-fonts \
      texlive-hmtrump-fonts \
      texlive-ibarra-fonts \
      texlive-ibygrk-fonts \
      texlive-imfellenglish-fonts \
      texlive-inconsolata-fonts \
      texlive-initials-fonts \
      texlive-inriafonts-fonts \
      texlive-inter-fonts \
      texlive-ipaex-fonts \
      texlive-ipaex-type1-fonts \
      texlive-iwona-fonts \
      texlive-jmn-fonts \
      texlive-josefin-fonts \
      texlive-junicode-fonts \
      texlive-kerkis-fonts \
      texlive-knitting-fonts \
      texlive-kpfonts-fonts \
      texlive-kpfonts-otf-fonts \
      texlive-kurier-fonts \
      texlive-latex-fonts \
      texlive-lato-fonts \
      texlive-lexend-fonts \
      texlive-libertine-fonts \
      texlive-libertinus-fonts \
      texlive-libertinus-fonts-fonts \
      texlive-libertinus-type1-fonts \
      texlive-libertinust1math-fonts \
      texlive-librebaskerville-fonts \
      texlive-librebodoni-fonts \
      texlive-librecaslon-fonts \
      texlive-librefranklin-fonts \
      texlive-libris-fonts \
      texlive-lilyglyphs-fonts \
      texlive-linearA-fonts \
      texlive-linguisticspro-fonts \
      texlive-lm-fonts \
      texlive-lm-math-fonts \
      texlive-lobster2-fonts \
      texlive-logix-fonts \
      texlive-lxfonts-fonts \
      texlive-magra-fonts \
      texlive-manfnt-font-fonts \
      texlive-marcellus-fonts \
      texlive-marvosym-fonts \
      texlive-mathabx-type1-fonts \
      texlive-mathdesign-fonts \
      texlive-mathpazo-fonts \
      texlive-mdsymbol-fonts \
      texlive-merriweather-fonts \
      texlive-metapost-fonts \
      texlive-mflogo-font-fonts \
      texlive-miama-fonts \
      texlive-mintspirit-fonts \
      texlive-missaali-fonts \
      texlive-mlmodern-fonts \
      texlive-mnsymbol-fonts \
      texlive-montex-fonts \
      texlive-montserrat-fonts \
      texlive-musixtex-fonts \
      texlive-musixtex-fonts-fonts \
      texlive-mxedruli-fonts \
      texlive-nanumtype1-fonts \
      texlive-ncntrsbk-fonts \
      texlive-newcomputermodern-fonts \
      texlive-newpx-fonts \
      texlive-newtx-fonts \
      texlive-newtxsf-fonts \
      texlive-newtxtt-fonts \
      texlive-niceframe-type1-fonts \
      texlive-nimbus15-fonts \
      texlive-noto-emoji-fonts \
      texlive-noto-fonts \
      texlive-notomath-fonts \
      texlive-novel-fonts \
      texlive-nunito-fonts \
      texlive-ocherokee-fonts \
      texlive-ocr-b-outline-fonts \
      texlive-oinuit-fonts \
      texlive-old-arrows-fonts \
      texlive-oldstandard-fonts \
      texlive-omega-fonts \
      texlive-opensans-fonts \
      texlive-oswald-fonts \
      texlive-overlock-fonts \
      texlive-padauk-fonts \
      texlive-palatino-fonts \
      texlive-paratype-fonts \
      texlive-pdftex-fonts \
      texlive-phaistos-fonts \
      texlive-philokalia-fonts \
      texlive-pigpen-fonts \
      texlive-pl-fonts \
      texlive-playfair-fonts \
      texlive-plex-fonts \
      texlive-plimsoll-fonts \
      texlive-poiretone-fonts \
      texlive-poltawski-fonts \
      texlive-prodint-fonts \
      texlive-ptex-fonts \
      texlive-punknova-fonts \
      texlive-pxfonts-fonts \
      texlive-qualitype-fonts \
      texlive-quattrocento-fonts \
      texlive-raleway-fonts \
      texlive-recycle-fonts \
      texlive-roboto-fonts \
      texlive-rojud-fonts \
      texlive-romande-fonts \
      texlive-rosario-fonts \
      texlive-rsfs-fonts \
      texlive-sanskrit-t1-fonts \
      texlive-sansmathfonts-fonts \
      texlive-scanpages-fonts \
      texlive-scholax-fonts \
      texlive-semaphor-fonts \
      texlive-shobhika-fonts \
      texlive-skaknew-fonts \
      texlive-sourcecodepro-fonts \
      texlive-sourcesanspro-fonts \
      texlive-sourceserifpro-fonts \
      texlive-spectral-fonts \
      texlive-starfont-fonts \
      texlive-staves-fonts \
      texlive-step-fonts \
      texlive-stepgreek-fonts \
      texlive-stickstoo-fonts \
      texlive-stix-fonts \
      texlive-stix2-otf-fonts \
      texlive-stix2-type1-fonts \
      texlive-stmaryrd-fonts \
      texlive-superiors-fonts \
      texlive-svrsymbols-fonts \
      texlive-symbol-fonts \
      texlive-tabvar-fonts \
      texlive-tapir-fonts \
      texlive-tempora-fonts \
      texlive-tex-gyre-fonts \
      texlive-tex-gyre-math-fonts \
      texlive-tfrupee-fonts \
      texlive-theanodidot-fonts \
      texlive-theanomodern-fonts \
      texlive-theanooldstyle-fonts \
      texlive-times-fonts \
      texlive-tinos-fonts \
      texlive-tipa-fonts \
      texlive-trajan-fonts \
      texlive-twemoji-colr-fonts \
      texlive-txfonts-fonts \
      texlive-txfontsb-fonts \
      texlive-txuprcal-fonts \
      texlive-typicons-fonts \
      texlive-uhc-fonts \
      texlive-umtypewriter-fonts \
      texlive-unfonts-core-fonts \
      texlive-unfonts-extra-fonts \
      texlive-universalis-fonts \
      texlive-uptex-fonts \
      texlive-utopia-fonts \
      texlive-velthuis-fonts \
      texlive-venturisadf-fonts \
      texlive-vntex-fonts \
      texlive-wadalab-fonts \
      texlive-wasy-type1-fonts \
      texlive-xcharter-fonts \
      texlive-xits-fonts \
      texlive-xypic-fonts \
      texlive-yfonts-t1-fonts \
      texlive-yhmath-fonts \
      texlive-yinit-otf-fonts \
      texlive-zapfchan-fonts \
      texlive-zapfding-fonts \
      thai-fonts \
      thessalonica-oldstandard-otf-fonts \
      thessalonica-oldstandard-ttf-fonts \
      thessalonica-tempora-lgc-otf-fonts \
      thessalonica-tempora-lgc-ttf-fonts \
      thessalonica-theano-otf-fonts \
      thessalonica-theano-ttf-fonts \
      thryomanes-fonts \
      tibetan-machine-uni-fonts \
      tiro-bangla-fonts \
      tiro-devahindi-fonts \
      tiro-devamarathi-fonts \
      tiro-devasanskrit-fonts \
      tiro-gurmukhi-fonts \
      tiro-indigo-fonts \
      tiro-kannada-fonts \
      tiro-tamil-fonts \
      tiro-telugu-fonts \
      tuladha-jejeg-fonts \
      tv-fonts \
      ubuntu-fonts \
      un-fonts \
      unifraktur-fonts \
      vlgothic-fonts \
      vollkorn-fonts \
      wang-fonts \
      wqy-bitmap-fonts \
      wqy-microhei-fonts \
      wqy-zenhei-fonts \
      x11-japanese-bitmap-fonts \
      xano-mincho-fonts \
      xorg-x11-fonts
    if ! is_package_pattern_installed "$package_pattern" $package_count; then
      printf_red "⛈️  Failed to install fonts. Aborting."
      exit 1
    else
      printf_bold "✅  Fonts installed successfully."
    fi
  else
    printf_bold "✅  Found ${package_count} fonts matching the pattern '${package_pattern}'. Skipping."
  fi
}

install_package "wget"
install_package "git"

install_package "python311"
install_pip_package "pipenv" "2023.6.12"

install_package "nodejs18"
install_package "npm18"
install_corepack

install_package "go1.20"
install_swag
install_golangci
install_air

install_postgres
install_redis
install_minio
install_meilisearch
install_mailhog

install_package "exiftool"
install_package "ffmpeg-4"
install_package "poppler-tools"
install_package "ghostscript"
install_package "ImageMagick"

sudo bash -c "ulimit -n 1048576"
install_brew_package "ocrmypdf" "14.3.0"

install_tesseract
install_libreoffice
install_fonts
