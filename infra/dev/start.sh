#!/bin/bash

BASE_DIR="/opt"

start_cockroach() {
    cd "${BASE_DIR}/cockroach"
    echo "🚀  Starting CockroachDB..."
    sudo ./cockroach start-single-node --insecure --listen-addr=0.0.0.0:26257 --http-addr=0.0.0.0:8080 --background
}

start_redis() {
    sudo systemctl start redis
}

start_minio() {
    cd "${BASE_DIR}/minio"
    echo "🚀  Starting MinIO..."
    local log_dir="/var/log/minio"
    sudo mkdir -p $log_dir
    MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve sudo sh -c "'nohup /usr/local/bin/minio server ${BASE_DIR}/minio/data --console-address \":9001\" > ${log_dir}/log.txt 2>&1 &'"
}

start_meilisearch() {
    cd "${BASE_DIR}/meilisearch"
    echo "🚀  Starting Meilisearch..."
    local log_dir="/var/log/meilisearch"
    sudo mkdir -p $log_dir
    sudo sh -c "'nohup ./meilisearch --http-addr=0.0.0.0:7700 > ${log_dir}/log.txt 2>&1 &'"
}

start_mailhog() {
    cd "${BASE_DIR}/mailhog"
    echo "🚀  Starting MailHog..."
    local log_dir="/var/log/mailhog"
    sudo mkdir -p $log_dir
    sudo sh -c "'nohup ./MailHog_linux_amd64 > ${log_dir}/log.txt 2>&1 &'"
}

start_cockroach
start_redis
start_minio
start_meilisearch
start_mailhog
