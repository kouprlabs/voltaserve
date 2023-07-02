#!/bin/bash

base_dir="/opt"
mkdir -p $base_dir

git_branch=$(git symbolic-ref --short HEAD 2>/dev/null || echo "main")

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

check_supported_system() {
  local cpe_name
  cpe_name=$(grep -oP '(?<=^CPE_NAME=).+' /etc/os-release)
  cpe_name="${cpe_name//\"/}"
  local pretty_name
  pretty_name=$(grep -oP '(?<=^PRETTY_NAME=").*"' /etc/os-release | tr -d '"')
  if [[ $cpe_name == "cpe:/o:redhat:enterprise_linux:9:"* ||
    "$cpe_name" == "cpe:/o:rocky:rocky:9:"* ||
    "$cpe_name" == "cpe:/o:almalinux:almalinux:9:"* ||
    "$cpe_name" == "cpe:/o:oracle:linux:9:"* ]]; then
    printf_bold "✅  Found supported operating system '$pretty_name'"
  else
    printf_red "⛈️  Operating system not supported: ${pretty_name}"
    exit 1
  fi
}

is_package_pattern_installed() {
  local package_pattern="$1"
  local package_count="$2"
  local installed_count
  installed_count=$(dnf list installed "$package_pattern" | awk 'NF==3 {gsub(/\..*/, "", $1); print $1}' | wc -l)
  if [[ "$installed_count" -eq "$package_count" ]]; then
    return 0
  else
    return 1
  fi
}

install_postgres() {
  local postgres_service="postgresql"
  local not_found='! systemctl list-unit-files | grep -q '"${postgres_service}.service"''
  if eval "$not_found"; then
    printf_bold "📦  Installing service '${postgres_service}'..."
    sudo dnf module -y enable postgresql:15
    sudo dnf module -y install postgresql:15/server
    sudo postgresql-setup --initdb
    sudo systemctl enable --now "$postgres_service"
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

install_minio() {
  local minio_pkg="minio"
  local not_found="! rpm -q $minio_pkg >/dev/null"
  if eval "$not_found"; then
    printf_bold "📦  Installing package '${minio_pkg}'..."
    local minio_rpm="minio-20230629051228.0.0.x86_64.rpm"
    sudo wget -c "https://dl.min.io/server/minio/release/linux-amd64/${minio_rpm}" -P $base_dir
    sudo dnf install -y "${base_dir}/${minio_rpm}"
    sudo rm -f "${base_dir}/${minio_rpm}"
    sudo mkdir -p "${base_dir}/minio"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install package '${minio_pkg}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${minio_pkg}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${minio_pkg}' package. Skipping."
  fi
}

install_redis() {
  local redis_service="redis"
  local not_found='! systemctl list-unit-files | grep -q '"${redis_service}.service"''
  if eval "$not_found"; then
    printf_bold "📦  Installing service '${redis_service}'..."
    sudo dnf install -y $redis_service
    sudo systemctl enable $redis_service
    sudo systemctl start $redis_service
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install service '${redis_service}'. Aborting."
      exit 1
    else
      printf_bold "✅  Service '${redis_service}' installed successfully."
    fi
  else
    printf_bold "✅  Found service '$redis_service'. Skipping."
  fi
}

install_meilisearch() {
  local meilisearch_bin="${base_dir}/meilisearch/meilisearch"
  local not_found="! (command -v $meilisearch_bin >/dev/null 2>&1 && $meilisearch_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "📦  Installing binary '${meilisearch_bin}'..."
    sudo mkdir -p "${base_dir}/meilisearch"
    cd "${base_dir}/meilisearch" || exit 1
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
  local mailhog_bin="${base_dir}/mailhog/MailHog_linux_amd64"
  local not_found="! (command -v $mailhog_bin >/dev/null 2>&1 && $mailhog_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "📦  Installing binary '${mailhog_bin}'..."
    sudo mkdir -p "${base_dir}/mailhog"
    sudo wget -c https://github.com/mailhog/MailHog/releases/download/v1.0.1/MailHog_linux_amd64 -P "${base_dir}/mailhog"
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

install_dnf_package() {
  local package_name="$1"
  local extra_args="$2"
  local not_found="! dnf list installed $package_name &>/dev/null"
  if eval "$not_found"; then
    printf_bold "📦  Installing package '${package_name}'..."
    if [[ -n "$extra_args" ]]; then
      sudo dnf install -y "$package_name" "$extra_args"
    else
      sudo dnf install -y "$package_name"
    fi
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install package '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Package '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found package '${package_name}'. Skipping."
  fi
}

download_tesseract_trained_data() {
  local tessdata_dir="/usr/share/tesseract/tessdata"
  local file_path="${tessdata_dir}/$1.traineddata"
  local url="https://github.com/tesseract-ocr/tessdata/raw/4.1.0/$1.traineddata"
  if [[ ! -f "$file_path" ]]; then
    printf_bold "🧠  Downloading Tesseract trained data '${file_path}'..."
    sudo wget -c "$url" -P $tessdata_dir
    if [[ ! -f "$file_path" ]]; then
      printf_red "⛈️  Failed to download Tesseract trained data '${file_path}'. Aborting."
      exit 1
    else
      printf_bold "✅  Tesseract trained data '${file_path}' downloaded successfully."
    fi
  else
    printf_bold "✅  Found Tesseract trained data '${file_path}'. Skipping."
  fi
}

install_rpm_repository() {
  local repository_name="$1"
  local url="$2"
  local not_found="! dnf repolist | grep -q $repository_name"
  if eval "$not_found"; then
    printf_bold "🪐  Installing repository '${repository_name}'..."
    sudo dnf install -y "$url"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install repository '${repository_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Repository '${repository_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found repository '${repository_name}'. Skipping."
  fi
}

install_code_ready_builder_repository() {
  local cpe_name
  cpe_name=$(grep -oP '(?<=^CPE_NAME=).+' /etc/os-release)
  cpe_name="${cpe_name//\"/}"
  local arch
  arch=$(uname -m)
  if [[ $cpe_name == "cpe:/o:redhat:enterprise_linux:9:"* ]]; then
    local repo="codeready-builder-for-rhel-9-${arch}-rpms"
    local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
    if eval "$not_found"; then
      printf_bold "🪐  Installing repository '${repo}'..."
      sudo dnf config-manager --set-enabled "codeready-builder-for-rhel-9-${arch}-rpms"
      if eval "$not_found"; then
        printf_red "⛈️  Failed to install repository '${repo}'. Aborting."
        exit 1
      else
        printf_bold "✅  Repository '${repo}' installed successfully."
      fi
    else
      printf_bold "✅  Found repository '$repo'. Skipping."
    fi
  elif [[ $cpe_name == "cpe:/o:rocky:rocky:9:"* || $cpe_name == "cpe:/o:almalinux:almalinux:9:"* ]]; then
    local repo="crb"
    local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
    if eval "$not_found"; then
      printf_bold "🪐  Installing repository '${repo}'..."
      sudo dnf config-manager --set-enabled crb
      if eval "$not_found"; then
        printf_red "⛈️  Failed to install repository '${repo}'. Aborting."
        exit 1
      else
        printf_bold "✅  Repository '${repo}' installed successfully."
      fi
    else
      printf_bold "✅  Found repository '$repo'. Skipping."
    fi
  elif [[ $cpe_name == "cpe:/o:oracle:linux:9:"* ]]; then
    local repo="ol9_codeready_builder"
    local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
    if eval "$not_found"; then
      printf_bold "🪐  Installing repository '${repo}'..."
      sudo dnf config-manager --set-enabled ol9_codeready_builder
      if eval "$not_found"; then
        printf_red "⛈️  Failed to install repository '${repo}'. Aborting."
        exit 1
      else
        printf_bold "✅  Repository '${repo}' installed successfully."
      fi
    else
      printf_bold "✅  Found repository '$repo'. Skipping."
    fi
  else
    printf_red "Cannot check or install EPEL repository. Unknown operating system detected."
  fi
}

install_jbig2enc() {
  local jbig2_bin="/usr/local/bin/jbig2"
  local not_found="! (command -v $jbig2_bin >/dev/null 2>&1 && $jbig2_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "🔨  Building binary '${jbig2_bin}'..."
    cd "$base_dir" || exit 1
    sudo git clone --branch 0.29 --single-branch https://github.com/kouprlabs/jbig2enc.git
    cd "${base_dir}/jbig2enc" || exit 1
    sudo ./autogen.sh
    sudo ./configure --with-extra-libraries=/usr/local/lib/ --with-extra-includes=/usr/local/include/
    sudo make
    sudo make install
    cd "$base_dir" || exit 1
    sudo rm -rf "${base_dir}/jbig2enc"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to build binary '${jbig2_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Binary '${jbig2_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found binary '${jbig2_bin}'. Skipping."
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
    printf_bold "✅  Found Python package '$package_name'. Skipping."
  fi
}

install_nodejs_18() {
  local not_found='! dnf list installed nodejs >/dev/null 2>&1 || ! node --version | grep -qE "^v18\."'
  if eval "$not_found"; then
    printf_bold "💎  Installing Node.js v18..."
    sudo dnf module -y enable nodejs:18
    sudo dnf module -y install nodejs:18/common
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install Node.js v18 '${package_name}'. Aborting."
      exit 1
    else
      printf_bold "✅  Node.js v18 '${package_name}' installed successfully."
    fi
  else
    printf_bold "✅  Found Node.js v18'. Skipping."
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
    cd "$HOME" || exit 1
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
    mkdir -p "${HOME}/bin"
    mv "$(go env GOPATH)/bin/swag" "${HOME}/bin/swag"
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
    mkdir -p "${HOME}/bin"
    mv "$(go env GOPATH)/bin/air" "${HOME}/bin/air"
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

install_fonts() {
  local package_pattern="*-fonts"
  local package_count=188
  if ! is_package_pattern_installed "$package_pattern" $package_count; then
    printf_bold "📦  Installing fonts..."
    sudo dnf install -y \
      google-noto-sans-cjk-ttc-fonts \
      google-noto-serif-cjk-ttc-fonts \
      aajohan-comfortaa-fonts \
      abattis-cantarell-fonts \
      adobe-source-code-pro-fonts \
      bitmap-fangsongti-fonts \
      culmus-aharoni-clm-fonts \
      culmus-caladings-clm-fonts \
      culmus-david-clm-fonts \
      culmus-drugulin-clm-fonts \
      culmus-ellinia-clm-fonts \
      culmus-frank-ruehl-clm-fonts \
      culmus-hadasim-clm-fonts \
      culmus-miriam-clm-fonts \
      culmus-miriam-mono-clm-fonts \
      culmus-nachlieli-clm-fonts \
      culmus-simple-clm-fonts \
      culmus-stamashkenaz-clm-fonts \
      culmus-stamsefarad-clm-fonts \
      culmus-yehuda-clm-fonts \
      dejavu-lgc-sans-fonts \
      dejavu-lgc-sans-mono-fonts \
      dejavu-lgc-serif-fonts \
      dejavu-sans-fonts \
      dejavu-sans-mono-fonts \
      dejavu-serif-fonts \
      fontawesome-fonts \
      ghostscript-tools-fonts \
      google-carlito-fonts \
      google-droid-sans-fonts \
      google-droid-sans-mono-fonts \
      google-droid-serif-fonts \
      google-noto-emoji-color-fonts \
      google-noto-emoji-fonts \
      google-noto-sans-armenian-fonts \
      google-noto-sans-avestan-fonts \
      google-noto-sans-bengali-fonts \
      google-noto-sans-bengali-ui-fonts \
      google-noto-sans-brahmi-fonts \
      google-noto-sans-carian-fonts \
      google-noto-sans-cherokee-fonts \
      google-noto-sans-coptic-fonts \
      google-noto-sans-deseret-fonts \
      google-noto-sans-devanagari-fonts \
      google-noto-sans-devanagari-ui-fonts \
      google-noto-sans-egyptian-hieroglyphs-fonts \
      google-noto-sans-ethiopic-fonts \
      google-noto-sans-fonts \
      google-noto-sans-georgian-fonts \
      google-noto-sans-glagolitic-fonts \
      google-noto-sans-gujarati-fonts \
      google-noto-sans-gujarati-ui-fonts \
      google-noto-sans-gurmukhi-fonts \
      google-noto-sans-hebrew-fonts \
      google-noto-sans-imperial-aramaic-fonts \
      google-noto-sans-kaithi-fonts \
      google-noto-sans-kannada-fonts \
      google-noto-sans-kannada-ui-fonts \
      google-noto-sans-kayah-li-fonts \
      google-noto-sans-kharoshthi-fonts \
      google-noto-sans-khmer-fonts \
      google-noto-sans-khmer-ui-fonts \
      google-noto-sans-lao-fonts \
      google-noto-sans-lao-ui-fonts \
      google-noto-sans-lycian-fonts \
      google-noto-sans-lydian-fonts \
      google-noto-sans-malayalam-fonts \
      google-noto-sans-malayalam-ui-fonts \
      google-noto-sans-mono-fonts \
      google-noto-sans-nko-fonts \
      google-noto-sans-old-south-arabian-fonts \
      google-noto-sans-old-turkic-fonts \
      google-noto-sans-osmanya-fonts \
      google-noto-sans-phoenician-fonts \
      google-noto-sans-shavian-fonts \
      google-noto-sans-sinhala-fonts \
      google-noto-sans-sinhala-vf-fonts \
      google-noto-sans-symbols-fonts \
      google-noto-sans-tamil-fonts \
      google-noto-sans-tamil-ui-fonts \
      google-noto-sans-telugu-fonts \
      google-noto-sans-telugu-ui-fonts \
      google-noto-sans-thaana-fonts \
      google-noto-sans-thai-fonts \
      google-noto-sans-thai-ui-fonts \
      google-noto-sans-ugaritic-fonts \
      google-noto-sans-vai-fonts \
      google-noto-serif-armenian-fonts \
      google-noto-serif-fonts \
      google-noto-serif-georgian-fonts \
      google-noto-serif-gurmukhi-vf-fonts \
      google-noto-serif-khmer-fonts \
      google-noto-serif-lao-fonts \
      google-noto-serif-sinhala-vf-fonts \
      google-noto-serif-thai-fonts \
      google-roboto-slab-fonts \
      gubbi-fonts \
      ht-caladea-fonts \
      jomolhari-fonts \
      julietaula-montserrat-fonts \
      kacst-art-fonts \
      kacst-book-fonts \
      kacst-decorative-fonts \
      kacst-digital-fonts \
      kacst-farsi-fonts \
      kacst-letter-fonts \
      kacst-naskh-fonts \
      kacst-office-fonts \
      kacst-one-fonts \
      kacst-pen-fonts \
      kacst-poster-fonts \
      kacst-qurn-fonts \
      kacst-screen-fonts \
      kacst-title-fonts \
      kacst-titlel-fonts \
      khmer-os-battambang-fonts \
      khmer-os-bokor-fonts \
      khmer-os-content-fonts \
      khmer-os-fasthand-fonts \
      khmer-os-freehand-fonts \
      khmer-os-handwritten-fonts \
      khmer-os-metal-chrieng-fonts \
      khmer-os-muol-fonts \
      khmer-os-muol-pali-fonts \
      khmer-os-siemreap-fonts \
      khmer-os-system-fonts \
      lato-fonts \
      liberation-fonts \
      liberation-mono-fonts \
      liberation-narrow-fonts \
      liberation-sans-fonts \
      liberation-serif-fonts \
      libreoffice-opensymbol-fonts \
      lklug-fonts \
      lohit-assamese-fonts \
      lohit-bengali-fonts \
      lohit-devanagari-fonts \
      lohit-gujarati-fonts \
      lohit-gurmukhi-fonts \
      lohit-kannada-fonts \
      lohit-marathi-fonts \
      lohit-odia-fonts \
      lohit-tamil-fonts \
      lohit-telugu-fonts \
      madan-fonts \
      navilu-fonts \
      open-sans-fonts \
      overpass-fonts \
      paktype-naqsh-fonts \
      paktype-naskh-basic-fonts \
      paktype-tehreer-fonts \
      pt-sans-fonts \
      redhat-display-fonts \
      redhat-mono-fonts \
      redhat-text-fonts \
      saab-fonts \
      sil-abyssinica-fonts \
      sil-nuosu-fonts \
      sil-padauk-fonts \
      sil-scheherazade-fonts \
      smc-meera-fonts \
      smc-rachana-fonts \
      stix-fonts \
      texlive-latex-fonts \
      thai-scalable-garuda-fonts \
      thai-scalable-kinnari-fonts \
      thai-scalable-loma-fonts \
      thai-scalable-norasi-fonts \
      thai-scalable-purisa-fonts \
      thai-scalable-sawasdee-fonts \
      thai-scalable-tlwgmono-fonts \
      thai-scalable-tlwgtypewriter-fonts \
      thai-scalable-tlwgtypist-fonts \
      thai-scalable-tlwgtypo-fonts \
      thai-scalable-umpush-fonts \
      thai-scalable-waree-fonts \
      ucs-miscfixed-fonts \
      urw-base35-bookman-fonts \
      urw-base35-c059-fonts \
      urw-base35-d050000l-fonts \
      urw-base35-fonts \
      urw-base35-gothic-fonts \
      urw-base35-nimbus-mono-ps-fonts \
      urw-base35-nimbus-roman-fonts \
      urw-base35-nimbus-sans-fonts \
      urw-base35-p052-fonts \
      urw-base35-standard-symbols-ps-fonts \
      urw-base35-z003-fonts
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

show_next_steps() {
  echo
  printf_bold "🎉 You are ready to develop Voltaserve!"
  echo

  echo "1) Start infrastructure services:"
  printf_cyan "curl -sSfL \"https://raw.githubusercontent.com/kouprlabs/voltaserve/${git_branch}/scripts/rhel/start.sh?t=\$(date +%s)\" | sh -s"
  echo

  echo "2) Open a terminal in each microservice's subfolder, then start each one in development mode:"
  echo

  printf_grey "cd ./api"
  printf_magenta "go run ."
  echo

  printf_grey "cd ./conversion"
  printf_magenta "go run ."
  echo

  printf_grey "cd ./idp"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"
  echo

  printf_grey "cd ./webdav"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"
  echo

  printf_grey "cd ./ui"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"
  echo

  printf_grey "cd ./language"
  printf_magenta "pipenv install"
  printf_magenta "pipenv shell"
  printf_magenta "FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug"
  echo

  echo "Alternatively, if this is a VM you can use Visual Studio Code's remote development as described here: "
  printf_underlined "https://code.visualstudio.com/docs/remote/remote-overview"
  echo "For this you can find the workspace file (voltaserve.code-workspace) in the repository's root."
  echo

  echo "3) Stop infrastructure services (if needed):"
  printf_cyan "curl -sSfL \"https://raw.githubusercontent.com/kouprlabs/voltaserve/${git_branch}/scripts/rhel/stop.sh?t=\$(date +%s)\" | sh -s"
}

check_supported_system

install_dnf_package "tar"
install_dnf_package "wget"
install_dnf_package "git"

install_dnf_package "clang"
install_dnf_package "python3"
install_dnf_package "python3-devel"
install_dnf_package "python3-pip"
install_pip_package "pipenv" "2023.6.12"

install_nodejs_18
install_corepack

install_dnf_package "golang"
install_swag
install_golangci
install_air

install_postgres
install_redis
install_minio
install_meilisearch
install_mailhog

install_rpm_repository "epel" "https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm"
install_rpm_repository "rpmfusion-free-updates" "https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm"
install_rpm_repository "rpmfusion-nonfree-updates" "https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-9.noarch.rpm"
install_dnf_package "dnf-plugins-core"
install_code_ready_builder_repository

install_dnf_package "perl-Image-ExifTool"
install_dnf_package "ffmpeg" "--allowerasing"
install_dnf_package "poppler-utils"
install_dnf_package "ghostscript"
install_dnf_package "ImageMagick"

install_dnf_package "tesseract"
download_tesseract_trained_data "osd"
download_tesseract_trained_data "eng"
download_tesseract_trained_data "deu"
download_tesseract_trained_data "fra"
download_tesseract_trained_data "nld"
download_tesseract_trained_data "ita"
download_tesseract_trained_data "spa"
download_tesseract_trained_data "por"
download_tesseract_trained_data "swe"
download_tesseract_trained_data "fin"
download_tesseract_trained_data "jpn"
download_tesseract_trained_data "chi_sim"
download_tesseract_trained_data "chi_tra"
download_tesseract_trained_data "hin"
download_tesseract_trained_data "rus"
download_tesseract_trained_data "ara"

install_dnf_package "unpaper"
install_dnf_package "pngquant"
install_pip_package "ocrmypdf" "14.3.0"

install_dnf_package "libreoffice"
install_fonts

show_next_steps
