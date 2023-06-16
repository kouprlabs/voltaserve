#!/bin/bash

BASE_DIR="/opt"

start_cockroach() {
    echo "🚀  Starting CockroachDB..."
    local opt_dir="${BASE_DIR}/cockroach"
    local log_dir="/var/log/cockroach"
    sudo mkdir -p $log_dir
    sudo sh -c ''"$opt_dir"'/cockroach start-single-node --insecure --listen-addr=0.0.0.0:26257 --http-addr=0.0.0.0:8080 --background > '"$log_dir"'/log.txt 2>&1 &'
}

start_redis() {
    sudo systemctl start redis
}

start_minio() {
    echo "🚀  Starting MinIO..."
    local opt_dir="${BASE_DIR}/minio"
    local log_dir="/var/log/minio"
    sudo mkdir -p $log_dir
    MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve sudo sh -c 'nohup /usr/local/bin/minio server '"$opt_dir"'/data --console-address ":9001" > '"$log_dir"'/log.txt 2>&1 &'
}

start_meilisearch() {
    echo "🚀  Starting Meilisearch..."
    local opt_dir="${BASE_DIR}/meilisearch"
    local log_dir="/var/log/meilisearch"
    sudo mkdir -p $log_dir
    sudo sh -c 'nohup '"$opt_dir"'/meilisearch --http-addr=0.0.0.0:7700 > '"$log_dir"'/log.txt 2>&1 &'
}

start_mailhog() {
    echo "🚀  Starting MailHog..."
    local opt_dir="${BASE_DIR}/mailhog"
    local log_dir="/var/log/mailhog"
    sudo mkdir -p $log_dir
    sudo sh -c 'nohup '"$opt_dir"'/MailHog_linux_amd64 > /var/log/mailhog/log.txt 2>&1 &'
}

start_cockroach
start_redis
start_minio
start_meilisearch
start_mailhog
