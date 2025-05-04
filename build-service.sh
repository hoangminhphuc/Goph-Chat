#!/bin/bash

# build-service.sh

set -e

NETWORK_NAME="your-network"                       # Docker network name
MYSQL_CONTAINER_NAME="your-mysql-container"      # MySQL container name
REDIS_CONTAINER_NAME="your-redis-container"      # Redis container name
APP_CONTAINER_NAME="your-app-container"          # Application container name

MYSQL_ROOT_PASSWORD="your_mysql_password"        # ⚠️ Set your MySQL root password
MYSQL_PORT="3307"                                # Port on host for MySQL
REDIS_PORT="6379"                                # Port on host for Redis
APP_PORT="8080"                                  # Port on host for your application

MYSQL_IMAGE="mysql:8"                            # MySQL Docker image
REDIS_IMAGE="bitnami/redis:latest"               # Redis Docker image
APP_IMAGE="your-app-image:latest"                # ⚠️ Replace with your built app image name


# Build Docker image for our application
echo "📦 Building application Docker image..."
docker build -t "$APP_IMAGE" -f Dockerfile .

# Create Docker network if not exists
if ! docker network ls | grep -q "$NETWORK_NAME"; then
  echo "Creating Docker network: $NETWORK_NAME"
  docker network create "$NETWORK_NAME"
else
  echo "Docker network $NETWORK_NAME already exists."
fi

# Start MySQL container
echo "Starting MySQL container..."
docker run -d \
  --name $MYSQL_CONTAINER_NAME \
  --network $NETWORK_NAME \
  -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
  -p $MYSQL_PORT:3306 \
  $MYSQL_IMAGE

# Start Redis container
echo "Starting Redis container..."
docker run -d \
  --name $REDIS_CONTAINER_NAME \
  --network $NETWORK_NAME \
  -e ALLOW_EMPTY_PASSWORD=yes \
  -p $REDIS_PORT:6379 \
  $REDIS_IMAGE

# Start Goph-Chat container
echo "Starting Goph-Chat container..."
docker run -d \
  --name $APP_CONTAINER_NAME \
  --network $NETWORK_NAME \
  -v "$(pwd)/.env:/app/.env" \
  --env-file .env \
  -e MYSQL_GORM_DB_URI="root:${MYSQL_ROOT_PASSWORD}@tcp(${MYSQL_CONTAINER_NAME}:3306)/goph-chat-db?charset=utf8mb4&parseTime=True&loc=Local" \
  -e REDIS_ADDR="${REDIS_CONTAINER_NAME}:6379" \
  -p $APP_PORT:8080 \
  $APP_IMAGE

echo "All services are up and running."
