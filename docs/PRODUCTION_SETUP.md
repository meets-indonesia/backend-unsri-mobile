# Production Setup Guide - Docker Compose

Panduan setup production dengan Docker Compose (Recommended untuk MVP).

## 🎯 Quick Answer

**Untuk MVP/Production awal: TIDAK perlu Kubernetes cluster.**

**Gunakan Docker Compose + Reverse Proxy (Nginx/Caddy).**

## ✅ Recommended Setup untuk Production

### Architecture

```
Internet
   │
   ▼
[Cloudflare/CDN] ← Optional, untuk DDoS protection
   │
   ▼
[Nginx/Caddy] ← Reverse Proxy, HTTPS, SSL
   │
   ▼
[Docker Compose] ← All services
   │
   ├──► API Gateway (:8080)
   ├──► Auth Service (:8081)
   ├──► User Service (:8082)
   └──► ... (other services)
   │
   ▼
[PostgreSQL + Redis + RabbitMQ]
```

## 🚀 Step-by-Step Setup

### Step 1: Prepare Server

```bash
# Update system
sudo apt-get update && sudo apt-get upgrade -y

# Install Docker (jika belum)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt-get install -y docker-compose-plugin

# Add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### Step 2: Clone & Setup Project

```bash
# Clone project (jika belum)
cd ~/sinergi/backend-unsri-mobile

# Create .env file untuk production
cat > .env << EOF
# Database
DATABASE_PASSWORD=your-strong-password-here
DATABASE_USER=unsri_user
DATABASE_NAME=unsri_db

# JWT
JWT_SECRET=$(openssl rand -base64 32)

# Redis
REDIS_PASSWORD=$(openssl rand -base64 32)

# RabbitMQ
RABBITMQ_PASSWORD=$(openssl rand -base64 32)

# Swagger (set false untuk production)
ENABLE_SWAGGER=false
EOF

# Set permissions
chmod 600 .env
```

### Step 3: Run Database Migrations

```bash
# Start database first
docker-compose -f deployments/docker-compose/docker-compose.yml up -d postgres

# Wait for database to be ready
sleep 10

# Run migrations
make migrate-up
```

### Step 4: Deploy Services

```bash
# Build and start all services
docker-compose -f deployments/docker-compose/docker-compose.yml up -d --build

# Check status
docker-compose -f deployments/docker-compose/docker-compose.yml ps

# View logs
docker-compose -f deployments/docker-compose/docker-compose.yml logs -f
```

### Step 5: Setup Reverse Proxy (Nginx)

**Install Nginx:**
```bash
sudo apt-get install -y nginx
```

**Create Nginx Config:**
```bash
sudo nano /etc/nginx/sites-available/unsri-backend
```

**Config:**
```nginx
server {
    listen 80;
    server_name api.unsri.ac.id;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.unsri.ac.id;

    # SSL certificates (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/api.unsri.ac.id/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.unsri.ac.id/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Proxy settings
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://localhost:8080/health;
        access_log off;
    }
}
```

**Enable site:**
```bash
sudo ln -s /etc/nginx/sites-available/unsri-backend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### Step 6: Setup SSL (Let's Encrypt)

```bash
# Install Certbot
sudo apt-get install -y certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d api.unsri.ac.id

# Auto-renewal (already configured by certbot)
sudo certbot renew --dry-run
```

### Step 7: Setup Firewall

```bash
# Allow SSH, HTTP, HTTPS
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
sudo ufw status
```

## 🔧 Alternative: Setup dengan Caddy (Auto SSL)

Caddy lebih simple karena auto SSL:

```bash
# Install Caddy
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy

# Create Caddyfile
sudo nano /etc/caddy/Caddyfile
```

**Caddyfile:**
```
api.unsri.ac.id {
    reverse_proxy localhost:8080
    
    # Security headers
    header {
        X-Frame-Options "SAMEORIGIN"
        X-Content-Type-Options "nosniff"
        X-XSS-Protection "1; mode=block"
    }
}
```

**Start Caddy:**
```bash
sudo systemctl enable caddy
sudo systemctl start caddy
sudo systemctl status caddy
```

## 📊 Monitoring

### Setup Basic Monitoring

```bash
# Install monitoring tools
docker-compose -f deployments/docker-compose/docker-compose.yml up -d

# Check service health
curl http://localhost:8080/health

# Monitor logs
docker-compose -f deployments/docker-compose/docker-compose.yml logs -f api-gateway
```

### Resource Monitoring

```bash
# Check resource usage
docker stats

# Check disk usage
df -h
docker system df
```

## 🔄 Maintenance

### Update Services

```bash
# Pull latest code
cd ~/sinergi/backend-unsri-mobile
git pull

# Rebuild and restart
docker-compose -f deployments/docker-compose/docker-compose.yml up -d --build

# Run migrations if needed
make migrate-up
```

### Backup Database

```bash
# Create backup script
cat > backup-db.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

docker exec unsri-postgres pg_dump -U unsri_user unsri_db > $BACKUP_DIR/backup_$DATE.sql
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
EOF

chmod +x backup-db.sh

# Add to crontab (daily backup at 2 AM)
(crontab -l 2>/dev/null; echo "0 2 * * * /path/to/backup-db.sh") | crontab -
```

## 🚨 When to Migrate to Kubernetes

Migrate ke Kubernetes ketika:

1. **Traffic Growth**
   - > 5,000 concurrent users
   - Need auto-scaling

2. **High Availability**
   - Cannot tolerate downtime
   - Need 99.9%+ uptime

3. **Team Growth**
   - Have DevOps team
   - Multiple teams

4. **Advanced Features**
   - Service mesh
   - Advanced monitoring
   - Canary deployments

## 📚 Related Documentation

- [Production Deployment Strategy](./PRODUCTION_DEPLOYMENT_STRATEGY.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Kubernetes Installation](./KUBERNETES_INSTALLATION.md)

---

**Kesimpulan: Untuk MVP, Docker Compose sudah cukup. Migrate ke Kubernetes ketika ada kebutuhan yang jelas.**

