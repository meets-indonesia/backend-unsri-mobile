# 🚀 Deployment Guide

Panduan lengkap untuk deployment backend UNSRI ke berbagai environment.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Docker Deployment](#docker-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Production Checklist](#production-checklist)
6. [Monitoring & Logging](#monitoring--logging)
7. [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Tools

- Docker 20.10+
- Docker Compose 2.0+
- Kubernetes 1.24+ (untuk K8s deployment)
- kubectl (untuk K8s deployment)
- PostgreSQL 15+ (jika tidak menggunakan Docker)
- Redis 7+ (jika tidak menggunakan Docker)
- golang-migrate (untuk database migrations)

### System Requirements

**Minimum:**
- CPU: 4 cores
- RAM: 8GB
- Storage: 50GB SSD

**Recommended (Production):**
- CPU: 8+ cores
- RAM: 16GB+
- Storage: 100GB+ SSD

## Environment Setup

### 1. Environment Variables

**Enable Swagger UI:**
```bash
# Set ENABLE_SWAGGER=true untuk enable Swagger UI
ENABLE_SWAGGER=true
```

Lihat [Swagger Deployment Guide](./SWAGGER_DEPLOYMENT.md) untuk detail lengkap.

Buat file `.env` untuk production:

```bash
# Database
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_USER=unsri_user
DATABASE_PASSWORD=<strong-password>
DATABASE_NAME=unsri_db
DATABASE_SSLMODE=require

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=<redis-password>

# RabbitMQ
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=unsri_user
RABBITMQ_PASSWORD=<rabbitmq-password>

# JWT
JWT_SECRET=<generate-strong-secret-key>

# Service Ports
PORT=8080
AUTH_SERVICE_PORT=8081
USER_SERVICE_PORT=8082
# ... (other service ports)

# Storage
STORAGE_TYPE=local
STORAGE_BASE_PATH=/storage
STORAGE_BASE_URL=https://api.unsri.ac.id/files
STORAGE_MAX_SIZE=10485760

# Logging
LOG_LEVEL=info

# API Gateway URLs (for service-to-service communication)
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082
# ... (other service URLs)
```

### 2. Generate JWT Secret

```bash
# Generate strong secret key
openssl rand -base64 32
```

### 3. Database Setup

```bash
# Create database
createdb unsri_db

# Run migrations
make migrate-up

# Or manually
migrate -path migrations -database "postgres://unsri_user:password@localhost:5432/unsri_db?sslmode=disable" up
```

## Docker Deployment

### 1. Build Images

```bash
# Build all images
docker-compose -f deployments/docker-compose/docker-compose.yml build

# Or build specific service
docker-compose -f deployments/docker-compose/docker-compose.yml build api-gateway
```

### 2. Start Services

```bash
# Start all services
docker-compose -f deployments/docker-compose/docker-compose.yml up -d

# Start specific services
docker-compose -f deployments/docker-compose/docker-compose.yml up -d postgres redis rabbitmq

# View logs
docker-compose -f deployments/docker-compose/docker-compose.yml logs -f

# View specific service logs
docker-compose -f deployments/docker-compose/docker-compose.yml logs -f api-gateway
```

### 3. Health Checks

```bash
# Check all services
docker-compose -f deployments/docker-compose/docker-compose.yml ps

# Test API Gateway
curl http://localhost:8080/health

# Test individual services
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # User Service
```

### 4. Stop Services

```bash
# Stop all services
docker-compose -f deployments/docker-compose/docker-compose.yml down

# Stop and remove volumes
docker-compose -f deployments/docker-compose/docker-compose.yml down -v
```

## Kubernetes Deployment

### 1. Setup Kubernetes Cluster

```bash
# Create namespace
kubectl create namespace unsri-backend

# Create secrets
kubectl create secret generic unsri-secrets \
  --from-literal=jwt-secret=<jwt-secret> \
  --from-literal=db-password=<db-password> \
  --from-literal=redis-password=<redis-password> \
  --namespace=unsri-backend
```

### 2. Deploy Infrastructure

```bash
# Deploy PostgreSQL
kubectl apply -f deployments/kubernetes/postgres.yaml

# Deploy Redis
kubectl apply -f deployments/kubernetes/redis.yaml

# Deploy RabbitMQ
kubectl apply -f deployments/kubernetes/rabbitmq.yaml
```

### 3. Deploy Services

```bash
# Deploy all services
kubectl apply -f deployments/kubernetes/services/

# Or deploy individually
kubectl apply -f deployments/kubernetes/services/api-gateway.yaml
kubectl apply -f deployments/kubernetes/services/auth-service.yaml
# ... (other services)
```

### 4. Deploy Ingress

```bash
# Deploy ingress controller (nginx)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Deploy ingress rules
kubectl apply -f deployments/kubernetes/ingress.yaml
```

### 5. Verify Deployment

```bash
# Check pods
kubectl get pods -n unsri-backend

# Check services
kubectl get svc -n unsri-backend

# Check ingress
kubectl get ingress -n unsri-backend

# View logs
kubectl logs -f deployment/api-gateway -n unsri-backend
```

## Production Checklist

### Security

- [ ] Change all default passwords
- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable SSL/TLS for database connections
- [ ] Configure firewall rules
- [ ] Enable rate limiting
- [ ] Setup CORS properly
- [ ] Use secrets management (Kubernetes Secrets, AWS Secrets Manager, etc.)
- [ ] Enable HTTPS for all services
- [ ] Regular security updates

### Database

- [ ] Setup database backups
- [ ] Configure connection pooling
- [ ] Enable query logging (for debugging)
- [ ] Setup read replicas (if needed)
- [ ] Configure automatic failover

### Monitoring

- [ ] Setup application monitoring (Prometheus, Grafana)
- [ ] Configure log aggregation (ELK, Loki)
- [ ] Setup alerting
- [ ] Monitor resource usage
- [ ] Setup uptime monitoring

### Performance

- [ ] Configure Redis caching
- [ ] Setup CDN for static files
- [ ] Enable gzip compression
- [ ] Configure load balancing
- [ ] Setup auto-scaling

### Backup & Recovery

- [ ] Database backup strategy
- [ ] File storage backup
- [ ] Disaster recovery plan
- [ ] Regular backup testing

## Monitoring & Logging

### Prometheus Metrics

Services expose metrics at `/metrics` endpoint:

```bash
curl http://localhost:8080/metrics
```

### Log Aggregation

All services use structured logging. Configure log aggregation:

**ELK Stack:**
```yaml
# docker-compose.yml
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.8.0
    # ... config
  
  logstash:
    image: docker.elastic.co/logstash/logstash:8.8.0
    # ... config
  
  kibana:
    image: docker.elastic.co/kibana/kibana:8.8.0
    # ... config
```

**Loki (Lightweight):**
```yaml
services:
  loki:
    image: grafana/loki:latest
    # ... config
  
  promtail:
    image: grafana/promtail:latest
    # ... config
```

### Health Checks

All services have health check endpoints:

```bash
# API Gateway
curl http://localhost:8080/health

# Individual services
curl http://localhost:8081/health  # Auth
curl http://localhost:8082/health  # User
# ... (other services)
```

## Troubleshooting

### Common Issues

**1. Database Connection Error**

```bash
# Check database is running
docker ps | grep postgres

# Check connection
psql -h localhost -U unsri_user -d unsri_db

# Check logs
docker logs unsri-postgres
```

**2. Service Not Starting**

```bash
# Check logs
docker logs <service-name>

# Check environment variables
docker exec <service-name> env

# Check port conflicts
lsof -i :8080
```

**3. Migration Errors**

```bash
# Check migration status
migrate -path migrations -database "postgres://..." version

# Rollback last migration
make migrate-down

# Force version (if needed)
migrate -path migrations -database "postgres://..." force <version>
```

**4. High Memory Usage**

```bash
# Check resource usage
docker stats

# Adjust limits in docker-compose.yml
services:
  api-gateway:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

### Debug Mode

Enable debug logging:

```bash
# Set log level
export LOG_LEVEL=debug

# Restart service
docker-compose restart <service-name>
```

### Performance Tuning

**Database:**
```sql
-- Check slow queries
SELECT * FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10;

-- Analyze tables
ANALYZE;
```

**Redis:**
```bash
# Check memory usage
redis-cli INFO memory

# Check slow commands
redis-cli SLOWLOG GET 10
```

## CI/CD Pipeline

### GitHub Actions Example

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build images
        run: docker-compose build
      
      - name: Run tests
        run: make test
      
      - name: Deploy to production
        run: |
          kubectl set image deployment/api-gateway api-gateway=...
          kubectl rollout status deployment/api-gateway
```

## Rollback Procedure

### Docker Compose

```bash
# Stop current version
docker-compose down

# Checkout previous version
git checkout <previous-commit>

# Rebuild and start
docker-compose up -d --build
```

### Kubernetes

```bash
# Rollback deployment
kubectl rollout undo deployment/api-gateway -n unsri-backend

# Check rollout status
kubectl rollout status deployment/api-gateway -n unsri-backend

# View rollout history
kubectl rollout history deployment/api-gateway -n unsri-backend
```

## Support

Untuk bantuan lebih lanjut:
- Dokumentasi: [docs/](./)
- Issues: GitHub Issues
- Email: support@unsri.ac.id

