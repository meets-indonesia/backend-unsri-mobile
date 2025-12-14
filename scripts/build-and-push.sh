#!/bin/bash

# Script untuk build dan push images ke Docker registry
# Usage: ./scripts/build-and-push.sh [registry] [version]
# Example: ./scripts/build-and-push.sh your-registry.com v1.0.0

set -e

REGISTRY="${1:-your-registry.com}"
VERSION="${2:-latest}"

echo "🚀 Building and pushing images to $REGISTRY..."
echo "Version: $VERSION"
echo ""

# Login to registry
echo "🔑 Logging in to registry..."
read -p "Press Enter after logging in (run: docker login $REGISTRY)..."

SERVICES=(
  "api-gateway"
  "auth-service"
  "user-service"
  "attendance-service"
  "schedule-service"
  "qr-service"
  "course-service"
  "broadcast-service"
  "notification-service"
  "calendar-service"
  "location-service"
  "access-service"
  "quick-actions-service"
  "file-storage-service"
  "search-service"
  "report-service"
)

for service in "${SERVICES[@]}"; do
  echo ""
  echo "📦 Building $service..."
  docker build -t $REGISTRY/unsri-$service:$VERSION -f deployments/docker/Dockerfile.$service .
  
  echo "⬆️  Pushing $service..."
  docker push $REGISTRY/unsri-$service:$VERSION
  
  echo "✅ $service done"
done

echo ""
echo "✅ All images built and pushed successfully!"
echo ""
echo "📝 Next steps:"
echo "   1. Update docker-compose.yml to use images from registry"
echo "   2. Pull images on server: docker-compose pull"
echo "   3. Deploy: docker-compose up -d"

