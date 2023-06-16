#!/bin/bash

BASE_DIR="/opt"

start_cockroach() {
    echo "🚀  Starting CockroachDB..."
    sudo /opt/cockroach/cockroach start-single-node --insecure --listen-addr=0.0.0.0:26257 --http-addr=0.0.0.0:8080 --background
}

start_redis() {
    sudo systemctl start redis
}

start_minio() {
    echo "🚀  Starting MinIO..."
    local log_dir="/var/log/minio"
    sudo mkdir -p $log_dir
    MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve sudo sh -c 'nohup /usr/local/bin/minio server /opt/minio/data --console-address ":9001" > /var/log/minio/log.txt 2>&1 &'
}

start_meilisearch() {
    echo "🚀  Starting Meilisearch..."
    local log_dir="/var/log/meilisearch"
    sudo mkdir -p $log_dir
    sudo sh -c 'nohup /opt/meilisearch/meilisearch --http-addr=0.0.0.0:7700 > /var/log/meilisearch/log.txt 2>&1 &'
}

start_mailhog() {
    echo "🚀  Starting MailHog..."
    local log_dir="/var/log/mailhog"
    sudo mkdir -p $log_dir
    sudo sh -c 'nohup /opt/mailhog/MailHog_linux_amd64 > /var/log/mailhog/log.txt 2>&1 &'
}

start_cockroach
start_redis
start_minio
start_meilisearch
start_mailhog
