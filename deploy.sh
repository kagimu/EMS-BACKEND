#!/bin/bash

echo "\n * Starting...\n * build the latest source code"
git checkout staging
echo "\n * Pulling code updates\n"
git pull

docker run -it -d --env-file ./.env --restart unless-stopped --net=host -v "$(pwd)"/fileRepo:/fileRepo:rw -p 2053:2053 wholesome-app