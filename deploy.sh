#!/bin/bash

echo "\n * Starting...\n * build the latest source code"
git checkout staging
echo "\n * Pulling code updates\n"
git pull
docker build -t wholesome-app .
docker-compose up -d