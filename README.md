# UNSRI Backend - Mobile App Backend

Backend microservices untuk aplikasi mobile UNSRI dengan arsitektur microservices menggunakan Go.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Services](#services)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Deployment](#deployment)
- [Testing](#testing)
- [Contributing](#contributing)

## 🏗️ Overview

Backend ini menggunakan arsitektur microservices dengan teknologi:
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Message Queue**: RabbitMQ (untuk scheduled broadcasts & async notifications)
- **Container**: Docker
- **Orchestration**: Kubernetes

## 🏛️ Architecture

```
┌─────────────┐
│   Client    │
│  (Mobile)   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         API Gateway (:8080)         │
│      (Routing & Load Balancing)     │
└──────┬──────────────────────────────┘
       │
       ├──► Auth Service (:8081)
       ├──► User Service (:8082)
       ├──► Attendance Service (:8084)
       ├──► Schedule Service (:8083)
       ├──► QR Service (:8085)
       ├──► Course Service (:8089)
       ├──► Broadcast Service (:8086)
       ├──► Notification Service (:8087)
       ├──► Calendar Service (:8088)
       ├──► Location Service (:8090)
       ├──► Access Service (:8091)
       ├──► Quick Actions Service (:8092)
       ├──► File Storage Service (:8093)
       ├──► Search Service (:8094)
       └──► Report Service (:8095)
       │
       ▼
┌─────────────────────────────────────┐
│  PostgreSQL  │  Redis  │  RabbitMQ  │
└─────────────────────────────────────┘
```

## 📦 Services

### Core Services

1. **API Gateway** (`:8080`) - Entry point untuk semua request
2. **Authentication Service** (`:8081`) - Autentikasi dan otorisasi dengan JWT
3. **User Management Service** (`:8082`) - Manajemen data pengguna
4. **Attendance Service** (`:8084`) - Manajemen kehadiran lengkap

### Academic Services

5. **Schedule Service** (`:8083`) - Manajemen jadwal kelas
6. **QR Code Service** (`:8085`) - Generate dan validasi QR code untuk absensi
7. **Course Management Service** (`:8089`) - Manajemen mata kuliah
8. **Academic Calendar Service** (`:8088`) - Kalender akademik

### Communication Services

9. **Broadcast Service** (`:8086`) - Broadcast pesan ke pengguna
10. **Notification Service** (`:8087`) - Notifikasi real-time

### Location & Access Services

11. **Location Service** (`:8090`) - Tap in/out dengan geofencing
12. **Access Control Service** (`:8091`) - Kontrol akses gate

### Additional Services

13. **Quick Actions Service** (`:8092`) - Quick actions berdasarkan role
14. **File Storage Service** (`:8093`) - Upload dan manajemen file
15. **Search Service** (`:8094`) - Pencarian data akademik
16. **Report Service** (`:8095`) - Generate laporan

## 👥 Roles

Sistem mendukung 3 role:
- **Mahasiswa** - Student
- **Dosen** - Lecturer  
- **Staff** - Staff member

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+ (jika run lokal tanpa Docker)
- Redis 7+ (jika run lokal tanpa Docker)
- golang-migrate (untuk database migrations)

### 3 Langkah Cepat

1. **Start Database & Redis:**
```bash
cd deployments/docker-compose
docker-compose up -d postgres redis rabbitmq
```

2. **Run Migrations:**
```bash
# Dari root project
make migrate-up
```

3. **Run Services:**
```bash
# Opsi A: Manual (buka terminal terpisah untuk setiap service)
make run-api-gateway
make run-auth-service
# ... (service lainnya)

# Opsi B: Docker Compose (semua services)
docker-compose -f deployments/docker-compose/docker-compose.yml up --build
```

### Verify Installation

```bash
# Test API Gateway
curl http://localhost:8080/health

# Test Auth Service
curl http://localhost:8081/health
```

Lihat [QUICK_START.md](./QUICK_START.md) untuk panduan lebih detail, atau [docs/LOCAL_DEVELOPMENT.md](./docs/LOCAL_DEVELOPMENT.md) untuk panduan lengkap development.

## 📚 API Documentation

### Swagger UI

Setelah service berjalan, akses Swagger UI di:
```
http://localhost:8080/swagger/index.html
```

Untuk enable Swagger:
```bash
export ENABLE_SWAGGER=true
make run-api-gateway
```

### Postman Collection

Import Postman collection untuk testing API:

1. Buka Postman
2. Import file: `postman/UNSRI_Backend_API.postman_collection.json`
3. Import environment: `postman/UNSRI_Backend_Environment.postman_environment.json`
4. Set `base_url` variable sesuai environment Anda

**Collection includes:**
- Authentication (Register, Login)
- User Management
- Attendance
- QR Code
- Location (Tap In/Out)
- Search
- Reports
- Dan lainnya

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register user baru
- `POST /api/v1/auth/login` - Login dan dapatkan JWT token

#### Users
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update profile

#### Attendance
- `POST /api/v1/attendance/scan` - Scan QR code untuk absensi
- `GET /api/v1/attendance/history` - Get attendance history

#### QR Code
- `POST /api/v1/qr/class/generate` - Generate QR untuk absensi kelas
- `POST /api/v1/qr/access/generate` - Generate QR untuk akses gate

Lihat [API Documentation](./docs/API.md) untuk dokumentasi lengkap semua endpoints.

## 💻 Development

### Project Structure

```
backend-unsri-mobile/
├── cmd/                    # Service entry points
│   ├── api-gateway/
│   ├── auth-service/
│   └── ...
├── internal/               # Internal packages
│   ├── api-gateway/
│   ├── auth/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── repository/
│   │   ├── service/
│   │   └── middleware/
│   └── shared/            # Shared packages
│       ├── database/
│       ├── errors/
│       ├── logger/
│       ├── models/
│       └── utils/
├── migrations/            # Database migrations
├── deployments/          # Deployment configs
│   ├── docker/
│   ├── docker-compose/
│   └── kubernetes/
├── docs/                 # Documentation
├── postman/             # Postman collections
└── scripts/             # Utility scripts
```

### Running Services Locally

**Option 1: Manual (Recommended for Development)**

```bash
# Terminal 1
make run-api-gateway

# Terminal 2
make run-auth-service

# Terminal 3
make run-user-service

# ... (service lainnya)
```

**Option 2: Script (macOS/Linux)**

```bash
./scripts/run-all-services.sh
```

**Option 3: Docker Compose**

```bash
docker-compose -f deployments/docker-compose/docker-compose.yml up
```

### Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-version
```

### Code Generation

```bash
# Generate Swagger docs
swag init -g cmd/api-gateway/main.go

# Format code
make fmt

# Run linter
make lint
```

## 🐳 Deployment

### Docker Deployment

```bash
# Build images
docker-compose -f deployments/docker-compose/docker-compose.yml build

# Start services
docker-compose -f deployments/docker-compose/docker-compose.yml up -d

# View logs
docker-compose -f deployments/docker-compose/docker-compose.yml logs -f
```

### Kubernetes Deployment

```bash
# Create namespace
kubectl create namespace unsri-backend

# Deploy services
kubectl apply -f deployments/kubernetes/

# Check status
kubectl get pods -n unsri-backend
```

Lihat [docs/DEPLOYMENT.md](./docs/DEPLOYMENT.md) untuk panduan deployment lengkap.

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific service tests
go test ./internal/auth/...
```

### API Testing

1. **Using Postman:**
   - Import collection dari `postman/UNSRI_Backend_API.postman_collection.json`
   - Set environment variables
   - Run requests

2. **Using cURL:**
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "role": "mahasiswa",
    "name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 📊 Database Schema

Database menggunakan PostgreSQL dengan schema yang expandable untuk sistem akademik lengkap.

**Core Tables:**
- `users` - Data pengguna (mahasiswa, dosen, staff)
- `attendances` - Data kehadiran
- `schedules` - Jadwal kelas
- `courses` - Mata kuliah
- `broadcasts` - Broadcast messages
- `notifications` - Notifikasi
- Dan lainnya...

Lihat [docs/DATABASE_SCHEMA.md](./docs/DATABASE_SCHEMA.md) untuk dokumentasi lengkap database schema.

## 🔧 Configuration

Configuration menggunakan environment variables. Lihat file `config.go` di setiap service untuk daftar lengkap environment variables.

### Common Environment Variables

```bash
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=unsri_user
DATABASE_PASSWORD=unsri_pass
DATABASE_NAME=unsri_db
DATABASE_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key-change-in-production

# Logging
LOG_LEVEL=info
```

## 🔐 Security

### Production Checklist

- [ ] Change all default passwords
- [ ] Use strong JWT secret (32+ characters)
- [ ] Enable SSL/TLS for database connections
- [ ] Configure firewall rules
- [ ] Enable rate limiting
- [ ] Setup CORS properly
- [ ] Use secrets management
- [ ] Enable HTTPS for all services
- [ ] Regular security updates

### Generate JWT Secret

```bash
openssl rand -base64 32
```

## 📈 Monitoring & Logging

### Health Checks

All services expose health check endpoints:

```bash
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
# ... (other services)
```

### Logging

All services use structured logging. Logs can be aggregated using:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Loki + Grafana
- Cloud logging services

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

[Your License Here]

## 👥 Contributors

- [Your Team Here]

## 📞 Support

- **Documentation**: [docs/](./docs/)
- **Issues**: [GitHub Issues](https://github.com/your-org/backend-unsri-mobile/issues)
- **Email**: support@unsri.ac.id

## 🔗 Useful Links

- [Quick Start Guide](./QUICK_START.md)
- [Deployment Guide](./docs/DEPLOYMENT.md)
- [Local Development Guide](./docs/LOCAL_DEVELOPMENT.md)
- [Services Status](./SERVICES_STATUS.md)
- [API Documentation](./docs/API.md)
- [Database Schema](./docs/DATABASE_SCHEMA.md)

---

**Made with ❤️ for UNSRI**
