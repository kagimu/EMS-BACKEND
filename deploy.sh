#!/bin/bash

echo "\n * Starting...\n * build the latest source code"
git checkout staging
echo "\n * Pulling code updates\n"
git pull
docker build --build-arg ENV_FILE=.env -t wholesome-app .
docker-compose up -d