#!/bin/bash

sudo mkdir -p /var/lib/local/reddit-clone-example || exit 1
cd /var/lib/local/reddit-clone-example
docker-compose up -d
sleep 2
docker-compose stop app
docker-compose up -d
