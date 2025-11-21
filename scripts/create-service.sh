#!/bin/bash

# Script to create a new service with standard structure
# Usage: ./scripts/create-service.sh <service-name> <port>

SERVICE_NAME=$1
PORT=$2

if [ -z "$SERVICE_NAME" ] || [ -z "$PORT" ]; then
    echo "Usage: $0 <service-name> <port>"
    echo "Example: $0 broadcast-service 8086"
    exit 1
fi

SERVICE_DIR="internal/${SERVICE_NAME}"
CMD_DIR="cmd/${SERVICE_NAME}"

echo "Creating service: $SERVICE_NAME on port $PORT"

# Create directories
mkdir -p ${SERVICE_DIR}/{config,handler,service,repository,middleware}
mkdir -p ${CMD_DIR}

# Copy template files from auth-service (you can customize these)
echo "Service structure created for $SERVICE_NAME"
echo "Next steps:"
echo "1. Copy and modify files from an existing service (e.g., auth-service)"
echo "2. Update config with port $PORT"
echo "3. Implement handlers, services, and repositories"
echo "4. Add to docker-compose.yml and kubernetes manifests"
echo "5. Update API Gateway routes"

