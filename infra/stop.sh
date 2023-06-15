#!/bin/bash

BASE_DIR="/opt"

stop_cockroach() {
    echo "🛑  Stopping CockroachDB..."
    pkill -f cockroach
}

stop_redis() {
    echo "🛑  Stopping Redis..."
    systemctl stop redis
}

stop_minio() {
    echo "🛑  Stopping MinIO..."
    pkill -f minio
}

stop_meilisearch() {
    echo "🛑  Stopping Meilisearch..."
    pkill -f meilisearch
}

stop_mailhog() {
    echo "🛑  Stopping MailHog..."
    pkill -f MailHog_linux_amd64
}

stop_cockroach
stop_redis
stop_minio
stop_meilisearch
stop_mailhog
