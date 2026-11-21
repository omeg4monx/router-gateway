#!/usr/bin/env bash
set -e

if docker container inspect omeg4mon-router >/dev/null 2>&1; then
    docker stop omeg4mon-router
    docker rm omeg4mon-router
fi
docker build -t omeg4mon-router .
docker run -d --name omeg4mon-router -p 20128:20128 --env-file .env -v "${OMEG4MON_ROUTER_VOLUME:-omeg4mon-router-data}:/app/data" omeg4mon-router