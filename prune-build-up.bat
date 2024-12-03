@echo off
echo Pruning Docker system, including unused images, containers, networks, and volumes...
docker system prune -a --volumes -f

echo Docker system prune complete.
echo Building Docker image with --no-cache...
docker-compose build --no-cache

echo Starting Docker containers with docker-compose up...
docker-compose up

echo Done.