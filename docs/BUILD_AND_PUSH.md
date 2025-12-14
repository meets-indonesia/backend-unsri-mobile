# Build & Push ke Docker Registry

Panduan build di local dan push ke Docker registry, lalu pull di server production.

## 🎯 Workflow

```
Local Machine (Build)
    ↓
Docker Registry (Push)
    ↓
Production Server (Pull & Deploy)
```

## 📋 Prerequisites

1. **Docker installed** di local machine ✅
2. **Docker Registry account** (Docker Hub, GitHub Container Registry, atau private registry)
3. **Access ke registry** dari production server

## 🚀 Step-by-Step

### Step 1: Build Images di Local

```bash
# Build semua images
docker-compose -f deployments/docker-compose/docker-compose.yml build

# Atau build satu per satu
docker build -t your-registry/unsri-api-gateway:latest -f deployments/docker/Dockerfile.api-gateway .
docker build -t your-registry/unsri-auth-service:latest -f deployments/docker/Dockerfile.auth-service .
# ... (other services)
```

### Step 2: Tag Images untuk Registry

```bash
# Tag dengan registry name
docker tag unsri-backend/api-gateway:latest your-registry/unsri-api-gateway:v1.0.0
docker tag unsri-backend/auth-service:latest your-registry/unsri-auth-service:v1.0.0
# ... (other services)

# Atau auto-tag saat build
docker build -t your-registry/unsri-api-gateway:v1.0.0 -f deployments/docker/Dockerfile.api-gateway .
```

### Step 3: Login ke Docker Registry

**Docker Hub:**
```bash
docker login
# Username: your-username
# Password: your-password
```

**GitHub Container Registry:**
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin
```

**Private Registry:**
```bash
docker login your-registry.com
```

### Step 4: Push Images ke Registry

```bash
# Push semua images
docker push your-registry/unsri-api-gateway:v1.0.0
docker push your-registry/unsri-auth-service:v1.0.0
docker push your-registry/unsri-user-service:v1.0.0
# ... (push semua services)
```

**Atau push semua sekaligus:**
```bash
docker images | grep unsri | awk '{print $1":"$2}' | xargs -I {} docker push {}
```

### Step 5: Update docker-compose.yml di Server

Edit `deployments/docker-compose/docker-compose.yml` untuk pull dari registry:

```yaml
api-gateway:
  image: your-registry/unsri-api-gateway:v1.0.0  # ⬅️ Gunakan image dari registry
  # build:  # ⬅️ Comment build section
  #   context: ../..
  #   dockerfile: deployments/docker/Dockerfile.api-gateway
  container_name: unsri-api-gateway
  # ... rest of config
```

### Step 6: Pull & Deploy di Server

```bash
# Login ke registry di server
docker login your-registry.com

# Pull images
docker pull your-registry/unsri-api-gateway:v1.0.0
docker pull your-registry/unsri-auth-service:v1.0.0
# ... (pull semua services)

# Atau pull via docker-compose
docker-compose -f deployments/docker-compose/docker-compose.yml pull

# Deploy
docker-compose -f deployments/docker-compose/docker-compose.yml up -d
```

## 📦 Docker Registry Options

### 1. Docker Hub (Free, Public/Private)

**Pros:**
- ✅ Free untuk public
- ✅ Mudah setup
- ✅ Well-known

**Cons:**
- ❌ Rate limit untuk free tier
- ❌ Public by default

**Setup:**
```bash
docker login
docker tag unsri-api-gateway:latest your-username/unsri-api-gateway:v1.0.0
docker push your-username/unsri-api-gateway:v1.0.0
```

### 2. GitHub Container Registry (Free, Private)

**Pros:**
- ✅ Free untuk private
- ✅ Integrated dengan GitHub
- ✅ No rate limit

**Cons:**
- ❌ Perlu GitHub token

**Setup:**
```bash
# Create GitHub token with write:packages permission
export GITHUB_TOKEN=your-token
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin

docker tag unsri-api-gateway:latest ghcr.io/your-username/unsri-api-gateway:v1.0.0
docker push ghcr.io/your-username/unsri-api-gateway:v1.0.0
```

### 3. Private Registry (Self-hosted)

**Setup Harbor atau Docker Registry:**
```bash
# Setup private registry
docker run -d -p 5000:5000 --restart=always --name registry registry:2

# Tag dan push
docker tag unsri-api-gateway:latest localhost:5000/unsri-api-gateway:v1.0.0
docker push localhost:5000/unsri-api-gateway:v1.0.0
```

## 🔧 Script untuk Build & Push

### Build & Push Script

```bash
#!/bin/bash
# scripts/build-and-push.sh

REGISTRY="your-registry.com"
VERSION="${1:-latest}"

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

echo "Building and pushing images..."

for service in "${SERVICES[@]}"; do
  echo "Building $service..."
  docker build -t $REGISTRY/unsri-$service:$VERSION -f deployments/docker/Dockerfile.$service .
  
  echo "Pushing $service..."
  docker push $REGISTRY/unsri-$service:$VERSION
  
  echo "✅ $service done"
done

echo "All done!"
```

### Update docker-compose.yml Script

Script untuk update docker-compose.yml menggunakan images dari registry:

```bash
#!/bin/bash
# scripts/update-compose-for-registry.sh

REGISTRY="${1:-your-registry.com}"
VERSION="${2:-latest}"

SERVICES=(
  "api-gateway"
  "auth-service"
  # ... all services
)

for service in "${SERVICES[@]}"; do
  sed -i "s|build:.*|image: $REGISTRY/unsri-$service:$VERSION|g" deployments/docker-compose/docker-compose.yml
  sed -i '/context:/d' deployments/docker-compose/docker-compose.yml
  sed -i '/dockerfile:/d' deployments/docker-compose/docker-compose.yml
done
```

## 📝 Workflow yang Disarankan

### Development:
1. Build di local
2. Test di local
3. Push ke registry
4. Deploy ke staging server

### Production:
1. Pull dari registry (no build needed!)
2. Deploy dengan docker-compose
3. Lebih cepat dan reliable

## ✅ Benefits

1. **Lebih cepat di server** - Tidak perlu build
2. **Konsisten** - Same image di semua environment
3. **Efisien** - Build sekali, deploy dimana saja
4. **Version control** - Tag images dengan version

## 🔒 Security

- Gunakan private registry untuk production
- Setup authentication
- Scan images untuk vulnerabilities
- Use specific version tags (bukan `latest`)

## 📚 Related

- [Docker Hub Documentation](https://docs.docker.com/docker-hub/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Private Docker Registry](https://docs.docker.com/registry/)

