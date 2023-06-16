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
    if ! (sudo command -v $cockroach_bin >/dev/null 2>&1 && sudo $cockroach_bin --version >/dev/null 2>&1); then
        echo "📦  Installing binary '${cockroach_bin}'..."
        cockroach_filename="cockroach-v23.1.3.linux-amd64"
        cockroach_tgz="${cockroach_filename}.tgz"
        sudo wget "https://binaries.cockroachdb.com/${cockroach_tgz}" -P $BASE_DIR
        sudo tar -xzf "${BASE_DIR}/${cockroach_tgz}" -C $BASE_DIR --transform="s/^${cockroach_filename}/cockroach/"
        sudo rm -f "${BASE_DIR}/${cockroach_tgz}"
    else
        echo "✅  Found binary '${cockroach_bin}'. Skipping."
    fi
}

install_minio() {
    local minio_pkg="minio"
    if ! rpm -q $minio_pkg >/dev/null; then
        echo "📦  Installing package '$minio_pkg'..."
        local minio_rpm="minio-20230609073212.0.0.x86_64.rpm"
        sudo wget "https://dl.min.io/server/minio/release/linux-amd64/archive/${minio_rpm}" -P $BASE_DIR
        sudo dnf install -y "${BASE_DIR}/${minio_rpm}"
        sudo rm -f "${BASE_DIR}/${minio_rpm}"
        sudo mkdir -p "${BASE_DIR}/minio"
    else
        echo "✅  Found package '${minio_pkg}' package. Skipping."
    fi
}

install_redis() {
    local redis_service="redis"
    if ! systemctl list-unit-files | grep -q "${redis_service}.service"; then
        echo "📦  Installing service '${redis_service}'..."
        sudo dnf install -y $redis_service
        sudo systemctl enable $redis_service
        sudo systemctl start $redis_service
    else
        echo "✅  Found service '$redis_service'. Skipping."
    fi
}

install_meilisearch() {
    local meilisearch_bin="${BASE_DIR}/meilisearch/meilisearch"
    if ! (sudo command -v $meilisearch_bin >/dev/null 2>&1 && sudo $meilisearch_bin --version >/dev/null 2>&1); then
        echo "📦  Installing binary '${meilisearch_bin}'..."
        sudo mkdir -p "${BASE_DIR}/meilisearch"
        cd "${BASE_DIR}/meilisearch"
        sudo curl -L https://install.meilisearch.com | sh
    else
        echo "✅  Found binary '${meilisearch_bin}'. Skipping."
    fi
}

install_mailhog() {
    local mailhog_bin="${BASE_DIR}/mailhog/MailHog_linux_amd64"
    if ! (sudo command -v $mailhog_bin >/dev/null 2>&1 && sudo $mailhog_bin --version >/dev/null 2>&1); then
        echo "📦  Installing binary '${mailhog_bin}'..."
        sudo mkdir -p "${BASE_DIR}/mailhog"
        sudo wget https://github.com/mailhog/MailHog/releases/download/v1.0.1/MailHog_linux_amd64 -P "${BASE_DIR}/mailhog"
        sudo chmod +x $mailhog_bin
    else
        echo "✅  Found binary '${mailhog_bin}'. Skipping."
    fi
}

install_dnf_package() {
    local package_name="$1"
    local extra_args="$2"
    if ! dnf list installed $package_name &>/dev/null; then
        echo "📦  Installing package '${package_name}'..."
        sudo dnf install -y $package_name $extra_args
    else
        echo "✅  Found package '${package_name}'. Skipping."
    fi
}

download_tesseract_trained_data() {
    local tessdata_dir="/usr/share/tesseract/tessdata"
    local file_path="${tessdata_dir}/$1.traineddata"
    local url="https://github.com/kouprlabs/tessdata/raw/main/$1.traineddata"
    if [ ! -f $file_path ]; then
        echo "🧠  Downloading Tesseract trained data '${file_path}'..."
        sudo wget $url -P $tessdata_dir
    else
        echo "✅  Found Tesseract trained data '${file_path}'. Skipping."
    fi
}

install_rpm_repository() {
    local repository_name="$1"
    local url="$2"
    if ! dnf repolist | grep -q $repository_name; then
        echo "🪐  Installing repository '${repository_name}'..."
        sudo dnf install -y $url
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
        if ! dnf repolist | grep -q "^${repo//\./\\.}"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled codeready-builder-for-rhel-9-${arch}-rpms
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    elif [[ $cpe_name == "cpe:/o:rocky:rocky:9:"* || $cpe_name == "cpe:/o:almalinux:almalinux:9:"* ]]; then
        local repo="crb"
        if ! dnf repolist | grep -q "^${repo//\./\\.}"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled crb
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    elif [[ $cpe_name == "cpe:/o:oracle:linux:9:"* ]]; then
        local repo="ol9_codeready_builder"
        if ! dnf repolist | grep -q "^${repo//\./\\.}"; then
            echo "🪐  Installing repository '${repo}'..."
            sudo dnf config-manager --set-enabled ol9_codeready_builder
        else
            echo "✅  Found repository '$repo'. Skipping."
        fi
    else
        echo "Cannot check or install EPEL repository. Unknown operating system detected."
    fi
}

install_jbig2enc() {
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
}

install_pip_package() {
    local package_name="$1"
    if ! pip show "$package_name" >/dev/null 2>&1; then
        echo "🐍  Installing Python package '${package_name}'..."
        pip3 install $package_name
    else
        echo "✅  Found Python package '$package_name'. Skipping."
    fi
}

install_nodejs_18() {
    if ! dnf list installed nodejs >/dev/null 2>&1 || ! node --version | grep -qE "^v18\."; then
        echo "💎  Installing Node.js v18..."
        sudo dnf module -y enable nodejs:18
        sudo dnf module -y install nodejs:18/common
    else
        echo "✅  Found Node.js v18'. Skipping."
    fi
}

install_corepack() {
    if ! npm list -g corepack >/dev/null 2>&1; then
        echo "💎  Installing NPM package 'corepack'..."
        sudo npm install -g corepack
    else
        echo "✅  Found NPM package 'corepack'. Skipping."
    fi
}

install_golangci() {
    local golangci_bin="$(go env GOPATH)/bin/golangci-lint"
    if ! (command -v $golangci_bin >/dev/null 2>&1 && $golangci_bin --version >/dev/null 2>&1); then
        echo "🐹  Installing Go binary '${golangci_bin}'..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.2
    else
        echo "✅  Found Go binary '${golangci_bin}'. Skipping."
    fi
}

install_swag() {
    local swag_bin="$(go env GOPATH)/bin/swag"
    if ! (command -v $swag_bin >/dev/null 2>&1 && $swag_bin --version >/dev/null 2>&1); then
        echo "🐹  Installing Go binary '${swag_bin}'..."
        go install github.com/swaggo/swag/cmd/swag@latest
        mkdir -p $HOME/bin
        mv $(go env GOPATH)/bin/swag $HOME/bin/swag
    else
        echo "✅  Found Go binary '${swag_bin}'. Skipping."
    fi
}

install_air() {
    local air_bin="$(go env GOPATH)/bin/air"
    if ! (command -v $air_bin >/dev/null 2>&1 && $air_bin -v >/dev/null 2>&1); then
        echo "🐹  Installing Go binary '${air_bin}'..."
        curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
    else
        echo "✅  Found Go binary '${air_bin}'. Skipping."
    fi
}

check_supported_system

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

install_pip_package "ocrmypdf"
install_pip_package "pipenv"

install_nodejs_18
install_corepack

install_air
install_golangci
install_swag

printf "\n🎉 \033[1mYou are ready to develop Voltaserve!\n\n\033[0m"

echo "1) Start infrastructure services:"
start_cmd="curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/start.sh?t=$(date +%s)" | sh -s"
printf "\033[32m${start_cmd}\n\n\033[0m"

echo "2) Create a user and database in CockroachDB (run only first time):"
user_and_db_cmd="curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/sql/create_user_and_database.sql?t=$(date +%s)" | /opt/cockroach/cockroach sql --insecure -u root"
printf "\033[36m${user_and_db_cmd}\n\n\033[0m"

echo "3) Create database schema (run only first time):"
schema_cmd="curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/sql/schema.sql?t=$(date +%s)" | /opt/cockroach/cockroach sql --insecure -u voltaserve"
printf "\033[36m${schema_cmd}\n\n\033[0m"

echo "4) Open a terminal in each microservice's subfolder, then start each one in development mode:"
echo

cd_cmd="cd ./api"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="go mod download && air"
printf "\033[35m${run_cmd}\n\n\033[0m"

cd_cmd="cd ./idp"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="pnpm i && pnpm dev"
printf "\033[35m${run_cmd}\n\n\033[0m"

cd_cmd="cd ./webdav"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="pnpm i && pnpm dev"
printf "\033[35m${run_cmd}\n\n\033[0m"

cd_cmd="cd ./conversion"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="go mod download && air"
printf "\033[35m${run_cmd}\n\n\033[0m"

cd_cmd="cd ./language"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="go mod download && air"
printf "\033[35m${run_cmd}\n\n\033[0m"

cd_cmd="cd ./ui"
printf "\033[90m${cd_cmd}\n\033[0m"
run_cmd="pnpm i && pnpm dev"
printf "\033[35m${run_cmd}\n\n\033[0m"

echo "5) To stop infrastructure services (if needed):"
stop_cmd="curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/stop.sh?t=$(date +%s)" | sh -s"
printf "\033[32m${stop_cmd}\n\n\033[0m"
