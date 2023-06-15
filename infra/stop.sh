#!/bin/bash

BASE_DIR="/opt"

stop_cockroach() {
    pkill -f cockroach
}

stop_redis() {
    systemctl stop redis

}

stop_minio() {
    pkill -f minio
}

stop_meilisearch() {
    pkill -f meilisearch
}

stop_mailhog() {
    pkill -f MailHog_linux_amd64
}

stop_cockroach
stop_redis
stop_minio
stop_meilisearch
stop_mailhog
