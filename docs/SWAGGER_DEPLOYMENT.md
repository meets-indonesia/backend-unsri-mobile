# Swagger Deployment Guide

Panduan untuk mengaktifkan Swagger UI di Docker dan Kubernetes deployment.

## ✅ Ya, Swagger Bisa Diakses!

Swagger UI **bisa diakses** di Docker dan Kubernetes jika environment variable `ENABLE_SWAGGER=true` di-set.

## 🐳 Docker Compose

### Enable Swagger

**Opsi 1: Set di .env file**
```bash
# Buat file .env di root project
echo "ENABLE_SWAGGER=true" >> .env
```

**Opsi 2: Set saat run**
```bash
ENABLE_SWAGGER=true docker-compose -f deployments/docker-compose/docker-compose.yml up -d
```

**Opsi 3: Edit docker-compose.yml langsung**
```yaml
api-gateway:
  environment:
    - ENABLE_SWAGGER=true  # Set ke true untuk enable
```

### Akses Swagger

Setelah container running:
```bash
# Local
http://localhost:8080/swagger/index.html

# Atau jika menggunakan domain
http://your-domain.com/swagger/index.html
```

### Verify

```bash
# Check environment variable
docker exec unsri-api-gateway env | grep ENABLE_SWAGGER

# Check logs
docker logs unsri-api-gateway | grep -i swagger

# Test endpoint
curl http://localhost:8080/swagger/index.html
```

## ☸️ Kubernetes

### Enable Swagger

**Opsi 1: Edit deployment file**

Edit `deployments/kubernetes/api-gateway.yaml`:
```yaml
env:
- name: ENABLE_SWAGGER
  value: "true"  # Set ke "false" untuk disable di production
```

**Opsi 2: Set via ConfigMap (Recommended)**

1. Buat ConfigMap:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: api-gateway-config
  namespace: unsri-backend
data:
  ENABLE_SWAGGER: "true"
  LOG_LEVEL: "info"
```

2. Update deployment untuk menggunakan ConfigMap:
```yaml
envFrom:
- configMapRef:
    name: api-gateway-config
```

**Opsi 3: Set via kubectl**
```bash
kubectl set env deployment/api-gateway ENABLE_SWAGGER=true -n unsri-backend
```

### Akses Swagger

**Via Ingress:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-gateway-ingress
  namespace: unsri-backend
spec:
  rules:
  - host: api.unsri.ac.id
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 80
```

Akses: `https://api.unsri.ac.id/swagger/index.html`

**Via Port Forward (Development):**
```bash
kubectl port-forward -n unsri-backend svc/api-gateway 8080:80
```

Akses: `http://localhost:8080/swagger/index.html`

### Verify

```bash
# Check environment variable
kubectl exec -n unsri-backend deployment/api-gateway -- env | grep ENABLE_SWAGGER

# Check logs
kubectl logs -n unsri-backend deployment/api-gateway | grep -i swagger

# Test endpoint
kubectl port-forward -n unsri-backend svc/api-gateway 8080:80
curl http://localhost:8080/swagger/index.html
```

## 🔒 Production Security

**⚠️ IMPORTANT:** Swagger UI menampilkan semua API endpoints dan bisa menjadi security risk di production.

### Recommended: Disable di Production

**Docker Compose:**
```yaml
environment:
  - ENABLE_SWAGGER=false  # Disable di production
```

**Kubernetes:**
```yaml
env:
- name: ENABLE_SWAGGER
  value: "false"  # Disable di production
```

### Alternative: Protect dengan Authentication

Tambahkan basic auth atau IP whitelist di Ingress:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-gateway-ingress
  annotations:
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: swagger-auth
    nginx.ingress.kubernetes.io/auth-realm: 'Swagger UI - Authentication Required'
spec:
  rules:
  - host: api.unsri.ac.id
    http:
      paths:
      - path: /swagger
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 80
```

## 📝 Generate Swagger Docs

Sebelum deploy, pastikan swagger docs sudah di-generate:

```bash
# Generate swagger docs
make swagger

# Atau manual
swag init -g cmd/api-gateway/main.go
```

Dockerfile sudah otomatis generate swagger docs saat build, tapi lebih baik generate dulu untuk memastikan.

## 🧪 Testing

### Test di Local Docker

```bash
# Enable swagger
export ENABLE_SWAGGER=true

# Start services
docker-compose -f deployments/docker-compose/docker-compose.yml up -d api-gateway

# Test
curl http://localhost:8080/swagger/index.html
```

### Test di Kubernetes

```bash
# Apply deployment
kubectl apply -f deployments/kubernetes/api-gateway.yaml

# Port forward
kubectl port-forward -n unsri-backend svc/api-gateway 8080:80

# Test
curl http://localhost:8080/swagger/index.html
```

## 🔍 Troubleshooting

### Swagger tidak muncul

1. **Check environment variable:**
```bash
# Docker
docker exec unsri-api-gateway env | grep ENABLE_SWAGGER

# Kubernetes
kubectl exec -n unsri-backend deployment/api-gateway -- env | grep ENABLE_SWAGGER
```

2. **Check logs:**
```bash
# Docker
docker logs unsri-api-gateway | grep -i swagger

# Kubernetes
kubectl logs -n unsri-backend deployment/api-gateway | grep -i swagger
```

3. **Check swagger docs exist:**
```bash
# Docker
docker exec unsri-api-gateway ls -la /root/docs/

# Kubernetes
kubectl exec -n unsri-backend deployment/api-gateway -- ls -la /root/docs/
```

4. **Regenerate swagger docs:**
```bash
make swagger
docker-compose -f deployments/docker-compose/docker-compose.yml build api-gateway
docker-compose -f deployments/docker-compose/docker-compose.yml up -d api-gateway
```

### 404 Not Found

Jika mendapat 404, pastikan:
- Environment variable `ENABLE_SWAGGER=true` sudah di-set
- Swagger docs sudah di-generate
- Route `/swagger/*any` sudah terdaftar di router

## 📚 Related Documentation

- [Swagger Setup Guide](./SWAGGER_SETUP.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [API Documentation](./API.md)

