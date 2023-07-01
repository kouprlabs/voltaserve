#!/bin/bash

base_dir="/opt"

start_postgres() {
    local not_found="! systemctl is-active --quiet redis"
    if eval "$not_found"; then
        echo "🚀  Starting PostgreSQL..."
        sudo systemctl start postgresql
        if eval "$not_found"; then
            echo "⛈️  Failed to start PostgreSQL."
        else
            echo "✅  PostgreSQL started successfully."
        fi
    else
        echo "✅  PostgreSQL is running. Skipping."
    fi
}

start_redis() {
    local not_found="! systemctl is-active --quiet redis"
    if eval "$not_found"; then
        echo "🚀  Starting Redis..."
        sudo systemctl start redis
        if eval "$not_found"; then
            echo "⛈️  Failed to start Redis."
        else
            echo "✅  Redis started successfully."
        fi
    else
        echo "✅  Redis is running. Skipping."
    fi
}

start_minio() {
    local not_found="! pgrep -f minio >/dev/null"
    if eval "$not_found"; then
        echo "🚀  Starting MinIO..."
        local opt_dir="${base_dir}/minio"
        local log_dir="/var/log/minio"
        sudo mkdir -p $log_dir
        cd $opt_dir || exit
        sudo sh -c 'MINIO_ROOT_USER=voltaserve MINIO_ROOT_PASSWORD=voltaserve nohup /usr/local/bin/minio server '"$opt_dir"'/data --console-address ":9001" > '"$log_dir"'/log.txt 2>&1 &'
        if eval "$not_found"; then
            echo "⛈️  Failed to start MinIO."
        else
            echo "✅  MinIO started successfully."
        fi
    else
        echo "✅  MinIO is running. Skipping."
    fi
}

start_meilisearch() {
    local not_found="! pgrep -f meilisearch >/dev/null"
    if eval "$not_found"; then
        echo "🚀  Starting Meilisearch..."
        local opt_dir="${base_dir}/meilisearch"
        local log_dir="/var/log/meilisearch"
        sudo mkdir -p $log_dir
        cd $opt_dir || exit
        sudo sh -c 'nohup '"$opt_dir"'/meilisearch --http-addr=0.0.0.0:7700 > '"$log_dir"'/log.txt 2>&1 &'
        if eval "$not_found"; then
            echo "⛈️  Failed to start Meilisearch."
        else
            echo "✅  Meilisearch started successfully."
        fi
    else
        echo "✅  Meilisearch is running. Skipping."
    fi
}

start_mailhog() {
    local not_found="! pgrep -f MailHog_linux_amd64 >/dev/null"
    if eval "$not_found"; then
        echo "🚀  Starting MailHog..."
        local opt_dir="${base_dir}/mailhog"
        local log_dir="/var/log/mailhog"
        sudo mkdir -p $log_dir
        cd $opt_dir || exit
        sudo sh -c 'nohup '"$opt_dir"'/MailHog_linux_amd64 > '"$log_dir"'/log.txt 2>&1 &'
        if eval "$not_found"; then
            echo "⛈️  Failed to start MailHog."
        else
            echo "✅  MailHog started successfully."
        fi
    else
        echo "✅  MailHog is running. Skipping."
    fi
}

start_postgres
start_redis
start_minio
start_meilisearch
start_mailhog
