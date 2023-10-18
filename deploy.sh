#!/bin/bash

echo "\n * Starting...\n * build the latest source code"
git checkout staging
echo "\n * Pulling code updates\n"
git pull

echo "Listing available running containers: "

docker ps
docker stop  wholesome_backend_cont

docker-compose -f go-app-compose.yaml  build

docker run -it -d --env-file ./.env --name=wholesome_backend_cont  --restart=unless-stopped::10  --net=host -v "$(pwd)"/fileRepo:/fileRepo:rw -p 2053:2053 wholesomeedbackend_wholesome-app