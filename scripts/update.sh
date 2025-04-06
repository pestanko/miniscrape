#! /bin/bash

# Stash all changes
git stash

# Pull the latest changes
git pull

# Build the docker compose
docker compose build

if [ $? -ne 0 ]; then
    echo "Failed to build docker compose"
    exit 1
fi

docker compose up -d --wait
if [ $? -ne 0 ]; then
    echo "Failed to start docker compose"
    exit 1
fi

echo "Docker compose started successfully"

# Pop the changes
git stash pop
