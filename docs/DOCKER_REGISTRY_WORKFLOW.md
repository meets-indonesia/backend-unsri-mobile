# Docker Registry Workflow

Panduan lengkap untuk build di local dan push ke Docker registry, lalu pull di server production.

## 🎯 Keuntungan

- ✅ **Lebih cepat di server** - Tidak perlu build (hemat waktu dan resource)
- ✅ **Konsisten** - Same image di semua environment
- ✅ **Efisien** - Build sekali, deploy dimana saja
- ✅ **Version control** - Tag images dengan version

## 📋 Workflow

```
Local Machine (Build Images)
    ↓
Docker Registry (Push Images)
    ↓
Production Server (Pull & Deploy)
```

## 🚀 Step-by-Step

### Step 1: Build di Local

```bash
# Build semua images
docker-compose -f deployments/docker-compose/docker-compose.yml build

# Atau gunakan script
chmod +x scripts/build-and-push.sh
./scripts/build-and-push.sh your-registry.com v1.0.0
```

### Step 2: Tag Images untuk Registry

```bash
# Tag dengan registry name
docker tag unsri-backend/api-gateway:latest your-registry.com/unsri-api-gateway:v1.0.0
docker tag unsri-backend/auth-service:latest your-registry.com/unsri-auth-service:v1.0.0
# ... (tag semua services)
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
export GITHUB_TOKEN=your-token
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin
```

**Private Registry:**
```bash
docker login your-registry.com
```

### Step 4: Push Images

```bash
# Push semua images
docker push your-registry.com/unsri-api-gateway:v1.0.0
docker push your-registry.com/unsri-auth-service:v1.0.0
# ... (push semua services)
```

### Step 5: Di Server - Pull & Deploy

```bash
# Login ke registry
docker login your-registry.com

# Pull images
docker-compose -f deployments/docker-compose/docker-compose.registry.yml pull

# Deploy
REGISTRY=your-registry.com VERSION=v1.0.0 docker-compose -f deployments/docker-compose/docker-compose.registry.yml up -d
```

## 📦 Docker Registry Options

### 1. Docker Hub (Recommended untuk Start)

```bash
# Login
docker login

# Tag & Push
docker tag unsri-api-gateway:latest your-username/unsri-api-gateway:v1.0.0
docker push your-username/unsri-api-gateway:v1.0.0
```

### 2. GitHub Container Registry (Recommended untuk Private)

```bash
# Create token di GitHub: Settings > Developer settings > Personal access tokens
# Permission: write:packages

# Login
echo $GITHUB_TOKEN | docker login ghcr.io -u your-username --password-stdin

# Tag & Push
docker tag unsri-api-gateway:latest ghcr.io/your-username/unsri-api-gateway:v1.0.0
docker push ghcr.io/your-username/unsri-api-gateway:v1.0.0
```

### 3. Private Registry (Self-hosted)

Setup Docker Registry sendiri di server.

## 🔧 Quick Commands

### Build & Push (Local)

```bash
# Login dulu
docker login your-registry.com

# Build dan push semua
REGISTRY=your-registry.com VERSION=v1.0.0 ./scripts/build-and-push.sh
```

### Pull & Deploy (Server)

```bash
# Set registry dan version
export REGISTRY=your-registry.com
export VERSION=v1.0.0

# Login
docker login $REGISTRY

# Pull images
docker-compose -f deployments/docker-compose/docker-compose.registry.yml pull

# Deploy
docker-compose -f deployments/docker-compose/docker-compose.registry.yml up -d
```

## ✅ Benefits

1. Server tidak perlu build (hemat CPU & RAM)
2. Build lebih cepat di local machine
3. Version control untuk images
4. Rollback mudah (pull version lama)

## 📝 Tips

- **Tag dengan version** - Jangan pakai `latest` untuk production
- **Private registry** - Untuk production, gunakan private registry
- **CI/CD** - Integrate dengan GitHub Actions untuk auto build & push

Lihat [Build & Push Guide](./BUILD_AND_PUSH.md) untuk detail lengkap.

