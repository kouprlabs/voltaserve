#!/bin/bash

BASE_DIR="/opt"
mkdir -p $BASE_DIR

check_supported_system() {
    local cpe_name=$(grep -oP '(?<=^CPE_NAME=).+' /etc/os-release)
    cpe_name="${cpe_name//\"/}"
    local pretty_name=$(grep -oP '(?<=^PRETTY_NAME=").*"' /etc/os-release | tr -d '"')
    if [[ $cpe_name == "cpe:/o:redhat:enterprise_linux:9:"* ||
        "$cpe_name" == "cpe:/o:rocky:rocky:9:"* ||
        "$cpe_name" == "cpe:/o:almalinux:almalinux:9:"* ||
        "$cpe_name" == "cpe:/o:oracle:linux:9:"* ]]; then
        echo "✅  Found supported operating system '$pretty_name'"
    else
        echo "⛈️  Operating system not supported: $pretty_name"
        exit 1
    fi
}

install_cockroach() {
    local cockroach_bin="${BASE_DIR}/cockroach/cockroach"
    local not_found="! (command -v $cockroach_bin >/dev/null 2>&1 && $cockroach_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "📦  Installing binary '${cockroach_bin}'..."
        cockroach_filename="cockroach-v23.1.3.linux-amd64"
        cockroach_tgz="${cockroach_filename}.tgz"
        sudo wget "https://binaries.cockroachdb.com/${cockroach_tgz}" -P $BASE_DIR
        sudo tar -xzf "${BASE_DIR}/${cockroach_tgz}" -C $BASE_DIR --transform="s/^${cockroach_filename}/cockroach/"
        sudo rm -f "${BASE_DIR}/${cockroach_tgz}"
        if eval "$not_found"; then
            echo "⛈️  Failed to install binary '${cockroach_bin}'. Aborting."
            exit 1
        else
            echo "✅  Binary '${cockroach_bin}' installed successfully."
        fi
    else
        echo "✅  Found binary '${cockroach_bin}'. Skipping."
    fi
}

install_minio() {
    local minio_pkg="minio"
    local not_found="! rpm -q $minio_pkg >/dev/null"
    if eval "$not_found"; then
        echo "📦  Installing package '$minio_pkg'..."
        local minio_rpm="minio-20230609073212.0.0.x86_64.rpm"
        sudo wget "https://dl.min.io/server/minio/release/linux-amd64/archive/${minio_rpm}" -P $BASE_DIR
        sudo dnf install -y "${BASE_DIR}/${minio_rpm}"
        sudo rm -f "${BASE_DIR}/${minio_rpm}"
        sudo mkdir -p "${BASE_DIR}/minio"
        if eval "$not_found"; then
            echo "⛈️  Failed to install package '${minio_pkg}'. Aborting."
            exit 1
        else
            echo "✅  Package '${minio_pkg}' installed successfully."
        fi
    else
        echo "✅  Found package '${minio_pkg}' package. Skipping."
    fi
}

install_redis() {
    local redis_service="redis"
    local not_found='! systemctl list-unit-files | grep -q '"${redis_service}.service"''
    if eval "$not_found"; then
        echo "📦  Installing service '${redis_service}'..."
        sudo dnf install -y $redis_service
        sudo systemctl enable $redis_service
        sudo systemctl start $redis_service
        if eval "$not_found"; then
            echo "⛈️  Failed to install service '${redis_service}'. Aborting."
            exit 1
        else
            echo "✅  Service '${redis_service}' installed successfully."
        fi
    else
        echo "✅  Found service '$redis_service'. Skipping."
    fi
}

install_meilisearch() {
    local meilisearch_bin="${BASE_DIR}/meilisearch/meilisearch"
    local not_found="! (command -v $meilisearch_bin >/dev/null 2>&1 && $meilisearch_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "📦  Installing binary '${meilisearch_bin}'..."
        sudo mkdir -p "${BASE_DIR}/meilisearch"
        cd "${BASE_DIR}/meilisearch"
        curl -L https://install.meilisearch.com | sudo sh
        sudo chmod +x $meilisearch_bin
        if eval "$not_found"; then
            echo "⛈️  Failed to install binary '${meilisearch_bin}'. Aborting."
            exit 1
        else
            echo "✅  Binary '${meilisearch_bin}' installed successfully."
        fi
    else
        echo "✅  Found binary '${meilisearch_bin}'. Skipping."
    fi
}

install_mailhog() {
    local mailhog_bin="${BASE_DIR}/mailhog/MailHog_linux_amd64"
    local not_found="! (command -v $mailhog_bin >/dev/null 2>&1 && $mailhog_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "📦  Installing binary '${mailhog_bin}'..."
        sudo mkdir -p "${BASE_DIR}/mailhog"
        sudo wget https://github.com/mailhog/MailHog/releases/download/v1.0.1/MailHog_linux_amd64 -P "${BASE_DIR}/mailhog"
        sudo chmod +x $mailhog_bin
        if eval "$not_found"; then
            echo "⛈️  Failed to install binary '${mailhog_bin}'. Aborting."
            exit 1
        else
            echo "✅  Binary '${mailhog_bin}' installed successfully."
        fi
    else
        echo "✅  Found binary '${mailhog_bin}'. Skipping."
    fi
}

install_dnf_package() {
    local package_name="$1"
    local extra_args="$2"
    local not_found="! dnf list installed $package_name &>/dev/null"
    if eval "$not_found"; then
        echo "📦  Installing package '${package_name}'..."
        sudo dnf install -y $package_name $extra_args
        if eval "$not_found"; then
            echo "⛈️  Failed to install package '${package_name}'. Aborting."
            exit 1
        else
            echo "✅  Package '${package_name}' installed successfully."
        fi
    else
        echo "✅  Found package '${package_name}'. Skipping."
    fi
}

download_tesseract_trained_data() {
    local tessdata_dir="/usr/share/tesseract/tessdata"
    local file_path="${tessdata_dir}/$1.traineddata"
    local url="https://github.com/kouprlabs/tessdata/raw/main/$1.traineddata"
    if [[ ! -f "$file_path" ]]; then
        echo "🧠  Downloading Tesseract trained data '${file_path}'..."
        sudo wget $url -P $tessdata_dir
        if [[ ! -f "$file_path" ]]; then
            echo "⛈️  Failed to download Tesseract trained data '${file_path}'. Aborting."
            exit 1
        else
            echo "✅  Tesseract trained data '${file_path}' downloaded successfully."
        fi
    else
        echo "✅  Found Tesseract trained data '${file_path}'. Skipping."
    fi
}

install_rpm_repository() {
    local repository_name="$1"
    local url="$2"
    local not_found="! dnf repolist | grep -q $repository_name"
    if eval "$not_found"; then
        echo "🪐  Installing repository '${repository_name}'..."
        sudo dnf install -y $url
        if eval "$not_found"; then
            echo "⛈️  Failed to install repository '${repository_name}'. Aborting."
            exit 1
        else
            echo "✅  Repository '${repository_name}' installed successfully."
        fi
    else
        echo "✅  Found repository '${repository_name}'. Skipping."
    fi
}

install_code_ready_builder_repository() {
    local cpe_name=$(grep -oP '(?<=^CPE_NAME=).+' /etc/os-release)
    cpe_name="${cpe_name//\"/}"
    local arch=$(uname -m)
    if [[ $cpe_name == "cpe:/o:redhat:enterprise_linux:9:"* ]]; then
        local repo="codeready-builder-for-rhel-9-${arch}-rpms"
        local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
        if eval "$not_found"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled codeready-builder-for-rhel-9-${arch}-rpms
            if eval "$not_found"; then
                echo "⛈️  Failed to install repository '${repo}'. Aborting."
                exit 1
            else
                echo "✅  Repository '${repo}' installed successfully."
            fi
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    elif [[ $cpe_name == "cpe:/o:rocky:rocky:9:"* || $cpe_name == "cpe:/o:almalinux:almalinux:9:"* ]]; then
        local repo="crb"
        local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
        if eval "$not_found"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled crb
            if eval "$not_found"; then
                echo "⛈️  Failed to install repository '${repo}'. Aborting."
                exit 1
            else
                echo "✅  Repository '${repo}' installed successfully."
            fi
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    elif [[ $cpe_name == "cpe:/o:oracle:linux:9:"* ]]; then
        local repo="ol9_codeready_builder"
        local not_found="! dnf repolist | grep -q "^${repo//\./\\.}""
        if eval "$not_found"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled ol9_codeready_builder
            if eval "$not_found"; then
                echo "⛈️  Failed to install repository '${repo}'. Aborting."
                exit 1
            else
                echo "✅  Repository '${repo}' installed successfully."
            fi
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    else
        echo "Cannot check or install EPEL repository. Unknown operating system detected."
    fi
}

install_jbig2enc() {
    local jbig2_bin="/usr/local/bin/jbig2"
    local not_found="! (command -v $jbig2_bin >/dev/null 2>&1 && $jbig2_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "🔨  Building binary '${jbig2_bin}'..."
        cd $BASE_DIR
        sudo git clone https://github.com/kouprlabs/jbig2enc.git
        cd "${BASE_DIR}/jbig2enc"
        sudo git checkout tags/0.29
        sudo ./autogen.sh
        sudo ./configure --with-extra-libraries=/usr/local/lib/ --with-extra-includes=/usr/local/include/
        sudo make
        sudo make install
        cd $BASE_DIR
        sudo rm -rf "${BASE_DIR}/jbig2enc"
        if eval "$not_found"; then
            echo "⛈️  Failed to build binary '${jbig2_bin}'. Aborting."
            exit 1
        else
            echo "✅  Binary '${jbig2_bin}' installed successfully."
        fi
    else
        echo "✅  Found binary '${jbig2_bin}'. Skipping."
    fi
}

install_pip_package() {
    local package_name="$1"
    local not_found="! pip show "$package_name" >/dev/null 2>&1"
    if eval "$not_found"; then
        echo "🐍  Installing Python package '${package_name}'..."
        pip3 install $package_name
        if eval "$not_found"; then
            echo "⛈️  Failed to install Python package '${package_name}'. Aborting."
            exit 1
        else
            echo "✅  Python package '${package_name}' installed successfully."
        fi
    else
        echo "✅  Found Python package '$package_name'. Skipping."
    fi
}

install_nodejs_18() {
    local not_found='! dnf list installed nodejs >/dev/null 2>&1 || ! node --version | grep -qE "^v18\."'
    if eval "$not_found"; then
        echo "💎  Installing Node.js v18..."
        sudo dnf module -y enable nodejs:18
        sudo dnf module -y install nodejs:18/common
        if eval "$not_found"; then
            echo "⛈️  Failed to install Node.js v18 '${package_name}'. Aborting."
            exit 1
        else
            echo "✅  Node.js v18 '${package_name}' installed successfully."
        fi
    else
        echo "✅  Found Node.js v18'. Skipping."
    fi
}

install_corepack() {
    local not_found="! npm list -g corepack >/dev/null 2>&1"
    if eval "$not_found"; then
        echo "💎  Installing NPM package 'corepack'..."
        sudo npm install -g corepack
        if eval "$not_found"; then
            echo "⛈️  Failed to install NPM package 'corepack'. Aborting."
            exit 1
        else
            echo "✅  NPM package 'corepack' installed successfully."
        fi
    else
        echo "✅  Found NPM package 'corepack'. Skipping."
    fi
}

install_golangci() {
    local golangci_bin="${HOME}/bin/golangci-lint"
    local not_found="! (command -v $golangci_bin >/dev/null 2>&1 && $golangci_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "🐹  Installing Go binary '${golangci_bin}'..."
        mkdir -p $HOME/bin
        cd $HOME
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.2
        if eval "$not_found"; then
            echo "⛈️  Failed to install Go binary '${golangci_bin}'. Aborting."
            exit 1
        else
            echo "✅  Go binary '${golangci_bin}' installed successfully."
        fi
    else
        echo "✅  Found Go binary '${golangci_bin}'. Skipping."
    fi
}

install_swag() {
    local swag_bin="${HOME}/bin/swag"
    local not_found="! (command -v $swag_bin >/dev/null 2>&1 && $swag_bin --version >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "🐹  Installing Go binary '${swag_bin}'..."
        go install github.com/swaggo/swag/cmd/swag@v1.8.12
        mkdir -p $HOME/bin
        mv $(go env GOPATH)/bin/swag $HOME/bin/swag
        if eval "$not_found"; then
            echo "⛈️  Failed to install Go binary '${swag_bin}'. Aborting."
            exit 1
        else
            echo "✅  Go binary '${swag_bin}' installed successfully."
        fi
    else
        echo "✅  Found Go binary '${swag_bin}'. Skipping."
    fi
}

install_air() {
    local air_bin="${HOME}/bin/air"
    local not_found="! (command -v $air_bin >/dev/null 2>&1 && $air_bin -v >/dev/null 2>&1)"
    if eval "$not_found"; then
        echo "🐹  Installing Go binary '${air_bin}'..."
        go install github.com/cosmtrek/air@v1.44.0
        mkdir -p $HOME/bin
        mv $(go env GOPATH)/bin/air $HOME/bin/air
        if eval "$not_found"; then
            echo "⛈️  Failed to install Go binary '${air_bin}'. Aborting."
            exit 1
        else
            echo "✅  Go binary '${air_bin}' installed successfully."
        fi
    else
        echo "✅  Found Go binary '${air_bin}'. Skipping."
    fi
}

printf_bold() {
    local msg="$1"
    printf "\033[1m${msg}\033[0m"
}

printf_cyan() {
    local msg="$1"
    printf "\033[36m${msg}\033[0m"
}

printf_grey() {
    local msg="$1"
    printf "\033[90m${msg}\033[0m"
}

printf_magenta() {
    local msg="$1"
    printf "\033[35m${msg}\033[0m"
}

show_next_steps() {
    printf_bold "\n🎉 You are ready to develop Voltaserve!\n\n"

    echo "1) Start infrastructure services:"
    local start_cmd='curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/start.sh?t=$(date +%s)" | sh -s'
    printf_cyan "${start_cmd}\n\n"

    echo "2) Create a user and database in CockroachDB (run only first time):"
    local user_and_db_cmd="curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/sql/create_user_and_database.sql?t=$(date +%s)" | /opt/cockroach/cockroach sql --insecure -u root"
    printf_cyan "${user_and_db_cmd}\n\n"

    echo "3) Create database schema (run only first time):"
    local schema_cmd='curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/sql/schema.sql?t=$(date +%s)" | /opt/cockroach/cockroach sql --insecure -u voltaserve'
    printf_cyan "${schema_cmd}\n\n"

    echo "4) Open a terminal in each microservice's subfolder, then start each one in development mode:"
    echo

    printf_grey "cd ./api\n"
    printf_magenta "go mod download && air\n\n"

    printf_grey "cd ./idp\n"
    printf_magenta "pnpm i && pnpm dev\n\n"

    printf_grey "cd ./webdav\n"
    printf_magenta "pnpm i && pnpm dev\n\n"

    printf_grey "cd ./conversion\n"
    printf_magenta "go mod download && air\n\n"

    printf_grey "cd ./language\n"
    printf_magenta "go mod download && air\n\n"

    printf_grey "cd ./ui\n"
    printf_magenta "pnpm i && pnpm dev\n\n"

    echo "5) To stop infrastructure services (if needed):"
    local stop_cmd='curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/stop.sh?t=$(date +%s)" | sh -s'
    printf_cyan "${stop_cmd}\n\n"
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

install_pip_package "ocrmypdf"
install_pip_package "pipenv"

install_nodejs_18
install_corepack

install_air
install_golangci
install_swag

show_next_steps
