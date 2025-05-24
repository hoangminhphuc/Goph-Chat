# Goph-Chat

A Go-based Chat Application that follows Clean Architecture principles, featuring user authentication, registration, and real-time messaging. It leverages Gin for routing, GORM for ORM with MySQL as the primary database, bcrypt for secure password hashing, JWT for token-based authentication, and Redis for caching and session management.

## Overview

This project implements a RESTful API built on Clean Architecture principles, which divides the code into distinct layers:

- **Business (Business Logic):** Core domain logic including registration, login, etc.
- **Repository (Data Access):** Interfaces with MySQL using GORM for persistent data operations.
- **Transport (Handlers):** Manages both RESTful HTTP routes and WebSocket connections using the Gin framework.

Additional features include:

- **Write-Behind Pattern:** Implements a Redis-backed write-behind pattern for efficient, deferred persistence to MySQL.
- **Two-Layer Caching:** Combines in-memory and Redis-based caching for faster data retrieval and reduced database load.
- **Pub/Sub Messaging:** Supports event-driven architecture for decoupled asynchronous processing.
- **Job Scheduling:** Define jobs and group them for execution in background routines.

## Getting Started

### Prerequisites

Ensure the following tools are installed:

1. **Go (>=1.18)**: For local development and builds
2. **Docker**: Required for containerization and running supporting services
3. **MySQL Workbench**: Used for easier data viewing and management (typically connected to MySQL running in Docker)

> **Note**: All supporting services (e.g., MySQL, Redis) can be run via Docker for a smoother setup.

### Installation

```bash
git clone https://github.com/hoangminhphuc/Goph-Chat.git
cd Goph-Chat
cp .env.example .env
go mod tidy
```

### Running Locally

```bash
go build && ./Task-Flow-System
```

Access the API at [http://localhost:8080](http://localhost:8080)

# Running with Docker

## Docker Compose

The easiest way to get this repository up and running is by using the `docker-compose.yaml` script.

> **Note:** Before running, make sure to update any placeholder variables in the file (e.g., container names, network name, ports, credentials) to match your system configuration.

### First-Time Setup (Build & Run)

```bash
docker-compose -f docker-compose.yaml up --build
# Add -d to run in detached mode:
# docker-compose -f docker-compose.yaml up --build -d
```
### Start Services (After Initial Setup)
```bash
docker-compose -f docker-compose.yaml start
```

## Automated Setup with `build-service.sh`

Instead of running each step manually, you can use the included `build-service.sh` script to build your Docker images and start the MySQL, Redis, and application containers with a single command. 

> **Note:** Be sure to open the script and update the placeholder variables just like the `docker-compose.yaml` above

### Run the automated setup

```bash
./build-service.sh
```
## Manual Setup with Docker

If you want more control over the Dockerization process, you can set up each service—MySQL, Redis, and the application—manually.

This approach is useful if you need to customize ports, volumes, environment variables, or container configurations beyond what's provided in the automated script

### 1. Create a Docker Network (if not already created):

```sh
docker network create chat-network
```

### 2. Start the MySQL Container:

```sh
docker run -d \
  --name <your-mysql-container-name> \  
  --network <your-network-name> \        # Must match the network used by the system container
  -e MYSQL_ROOT_PASSWORD=<your-password> \
  -p 3307:3306 \
  mysql:8
```

### 3. Start the Redis Container:

```sh
docker run -d \
  --name <your-redis-container-name> \
  --network <your-network-name> \        # Must match the network used by the system container
  -e ALLOW_EMPTY_PASSWORD=yes \
  -p 6379:6379 \
  bitnami/redis:latest
```

## Building the Docker Image

From the project root directory (where your `Dockerfile` and `.env` live), build your Go chat application image:

```sh
docker build -t goph-chat:1.0.0 -f Dockerfile .
```

## Single‑Container Deployment

Use the commands below to run your Go chat application container connected to both MySQL and Redis. Mount your `.env` file into the container so it picks up any additional configuration.

### On Linux:

```sh
docker run -d \
  --name <your-container-name> \
  --network <your-network-name> \
  -v $(pwd)/.env:/app/.env \
  --env-file .env \
  -e MYSQL_GORM_DB_URI="your-username:your-password@tcp(your-mysql-host:3306)/your-database?charset=utf8mb4&parseTime=True&loc=Local" \
  -e REDIS_ADDR="<your-redis-container-name>:6379" \
  -p 8080:8080 \
  goph-chat:1.0.0
```

### On Windows:

```sh
docker run -d \
  --name <your-container-name> \
  --network <your-network-name> \
  -v "path\to\env\file\.env:/app/.env" \
  --env-file .env \
  -e MYSQL_GORM_DB_URI="your-username:your-password@tcp(your-mysql-host:3306)/your-database?charset=utf8mb4&parseTime=True&loc=Local" \
  -e REDIS_ADDR="<your-redis-container-name>:6379" \
  -p 8080:8080 \
  goph-chat:1.0.0
```


Your system should now be live via Docker.
