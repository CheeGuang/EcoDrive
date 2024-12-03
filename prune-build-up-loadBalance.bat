@echo off
echo Pruning Docker system, including unused images, containers, networks, and volumes...
docker system prune -a --volumes -f

echo Docker system prune complete.
echo Building Docker images with --no-cache...
docker-compose build --no-cache

echo Scaling frontend service for load balancing...
docker-compose up --scale frontend=3

echo Starting the load balancer and other services...
docker-compose up load_balancer database authentication payment user vehicle

echo Docker containers are running with load balancing enabled.
echo Done.
