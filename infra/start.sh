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
    MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve sudo nohup /usr/local/bin/minio server "${BASE_DIR}/minio/data" --console-address ":9001" >log.txt &
}

start_meilisearch() {
    cd "${BASE_DIR}/meilisearch"
    echo "🚀  Starting Meilisearch..."
    sudo nohup ./meilisearch --http-addr=0.0.0.0:7700 >log.txt &
}

start_mailhog() {
    cd "${BASE_DIR}/mailhog"
    echo "🚀  Starting MailHog..."
    sudo nohup ./MailHog_linux_amd64 >log.txt &
}

start_cockroach
start_redis
start_minio
start_meilisearch
start_mailhog
