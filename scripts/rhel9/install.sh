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

install_cockroach() {
  local cockroach_bin="${base_dir}/cockroach/cockroach"
  local not_found="! (command -v $cockroach_bin >/dev/null 2>&1 && $cockroach_bin --version >/dev/null 2>&1)"
  if eval "$not_found"; then
    printf_bold "📦  Installing binary '${cockroach_bin}'..."
    cockroach_filename="cockroach-v23.1.3.linux-amd64"
    cockroach_tgz="${cockroach_filename}.tgz"
    sudo wget -c "https://binaries.cockroachdb.com/${cockroach_tgz}" -P $base_dir
    sudo tar -xzf "${base_dir}/${cockroach_tgz}" -C $base_dir --transform="s/^${cockroach_filename}/cockroach/"
    sudo rm -f "${base_dir}/${cockroach_tgz}"
    if eval "$not_found"; then
      printf_red "⛈️  Failed to install binary '${cockroach_bin}'. Aborting."
      exit 1
    else
      printf_bold "✅  Binary '${cockroach_bin}' installed successfully."
    fi
  else
    printf_bold "✅  Found binary '${cockroach_bin}'. Skipping."
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
    cd "${base_dir}/meilisearch" || exit
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
    sudo dnf install -y "$package_name" "$extra_args"
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
  local url="https://github.com/kouprlabs/tessdata/raw/4.1.0/$1.traineddata"
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
    cd "$base_dir" || exit
    sudo git clone --branch 0.29 --single-branch https://github.com/kouprlabs/jbig2enc.git
    cd "${base_dir}/jbig2enc" || exit
    sudo ./autogen.sh
    sudo ./configure --with-extra-libraries=/usr/local/lib/ --with-extra-includes=/usr/local/include/
    sudo make
    sudo make install
    cd "$base_dir" || exit
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
  printf_bold "🎉 You are ready to develop Voltaserve!"

  echo "1) Start infrastructure services:"
  printf_cyan "curl -sSfL \"https://raw.githubusercontent.com/kouprlabs/voltaserve/${git_branch}/scripts/rhel9/start.sh?t=$(date +%s)\" | sh -s"

  echo "2) Open a terminal in each microservice's subfolder, then start each one in development mode:"
  echo

  printf_grey "cd ./api"
  printf_magenta "air"

  printf_grey "cd ./conversion"
  printf_magenta "air"

  printf_grey "cd ./idp"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"

  printf_grey "cd ./webdav"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"

  printf_grey "cd ./ui"
  printf_magenta "pnpm install"
  printf_magenta "pnpm dev"

  printf_grey "cd ./language"
  printf_magenta "pipenv install"
  printf_magenta "pipenv shell"
  printf_magenta "FLASK_APP=server.py flask run --host=0.0.0.0 --port=5002 --debug"

  echo "Alternatively, if this is a VM you can use Visual Studio Code's remote development as described here: "
  printf_underlined "https://code.visualstudio.com/docs/remote/remote-overview"
  echo "For this you can find the workspace file (voltaserve.code-workspace) in the repository's root."

  echo "3) To stop infrastructure services (if needed):"
  printf_cyan "curl -sSfL \"https://raw.githubusercontent.com/kouprlabs/voltaserve/${git_branch}/scripts/rhel9/stop.sh?t=$(date +%s)\" | sh -s"
}

check_supported_system

install_dnf_package "tar"
install_dnf_package "wget"
install_dnf_package "git"

install_cockroach
install_meilisearch
install_mailhog
install_minio
install_redis

install_dnf_package "golang"
install_dnf_package "poppler-utils"
install_dnf_package "libreoffice"
install_dnf_package "python3-pip"
install_dnf_package "python3-devel"
install_dnf_package "ghostscript"
install_dnf_package "tesseract"
install_dnf_package "postgresql"

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
download_tesseract_trained_data "kor"
download_tesseract_trained_data "hin"
download_tesseract_trained_data "rus"
download_tesseract_trained_data "ara"

install_rpm_repository "epel" "https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm"
install_rpm_repository "rpmfusion-free-updates" "https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-9.noarch.rpm"
install_rpm_repository "rpmfusion-nonfree-updates" "https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-9.noarch.rpm"

install_dnf_package "GraphicsMagick"
install_dnf_package "pngquant"
install_dnf_package "unpaper"
install_dnf_package "perl-Image-ExifTool"
install_dnf_package "dnf-plugins-core"

install_code_ready_builder_repository

install_dnf_package "ffmpeg" "--allowerasing"
install_dnf_package "automake"
install_dnf_package "make"
install_dnf_package "autoconf"
install_dnf_package "libtool"
install_dnf_package "clang"
install_dnf_package "zlib"
install_dnf_package "zlib-devel"
install_dnf_package "libjpeg-turbo"
install_dnf_package "libjpeg-turbo-devel"
install_dnf_package "libwebp"
install_dnf_package "libwebp-devel"
install_dnf_package "libtiff"
install_dnf_package "libtiff-devel"
install_dnf_package "libpng"
install_dnf_package "libpng-devel"
install_dnf_package "leptonica-devel"

install_jbig2enc

install_pip_package "ocrmypdf" "14.2.1"
install_pip_package "pipenv" "2023.6.12"

install_nodejs_18
install_corepack

install_air
install_golangci
install_swag

show_next_steps
