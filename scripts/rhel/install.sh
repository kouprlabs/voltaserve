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

show_next_steps
