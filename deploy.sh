#!/bin/bash

echo "Listing available running containers: "

git pull

docker ps
docker stop  wholesome_backend_cont
docker rm wholesome_backend_cont
docker system prune -f

docker-compose -f go-app-compose.yaml  build

docker run -it -d --env-file ./.env --name=wholesome_backend_cont  --restart=unless-stopped  --net=host -v "$(pwd)"/fileRepo:/fileRepo:rw -p 2053:2053 wholesomeedbackend_wholesome-app

docker ps