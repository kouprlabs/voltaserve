#!/bin/bash

stop_postgres() {
    echo "🛑  Stopping PostgreSQL..."
    sudo systemctl stop postgresql
}

stop_redis() {
    echo "🛑  Stopping Redis..."
    sudo systemctl stop redis
}

stop_minio() {
    echo "🛑  Stopping MinIO..."
    sudo pkill -f minio
}

stop_meilisearch() {
    echo "🛑  Stopping Meilisearch..."
    sudo pkill -f meilisearch
}

stop_mailhog() {
    echo "🛑  Stopping MailHog..."
    sudo pkill -f MailHog_linux_amd64
}

stop_postgres
stop_redis
stop_minio
stop_meilisearch
stop_mailhog
