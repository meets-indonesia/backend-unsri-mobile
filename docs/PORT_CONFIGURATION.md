# Port Configuration

Dokumentasi konfigurasi port untuk UNSRI Backend.

## 📌 Port Mapping

### External Ports (Host)
- **PostgreSQL**: `5433` → Container port `5432`
- **Redis**: `6379` → Container port `6379`
- **RabbitMQ**: `5672` (AMQP), `15672` (Management UI)
- **API Gateway**: `8080`
- **Services**: `8081-8095`

### Internal Ports (Docker Network)
Services di dalam Docker network menggunakan:
- **PostgreSQL**: `5432` (via service name `postgres`)
- **Redis**: `6379` (via service name `redis`)
- **RabbitMQ**: `5672` (via service name `rabbitmq`)

## 🔧 Konfigurasi

### Untuk Services (di Docker)
Services connect via Docker network, jadi tetap pakai port 5432:
```yaml
DATABASE_HOST=postgres
DATABASE_PORT=5432  # ✅ Correct - via Docker network
```

### Untuk Migrations (dari Host)
Migrations run dari host, jadi pakai port 5433:
```bash
# Makefile sudah di-update
make migrate-up  # ✅ Otomatis pakai port 5433
```

### Manual Connection dari Host
```bash
# Connect ke PostgreSQL dari host
psql -h localhost -p 5433 -U unsri_user -d unsri_db

# Atau via Docker
docker exec -it unsri-postgres psql -U unsri_user -d unsri_db
```

## ⚠️ Catatan

- **Port 5432** di host sudah digunakan oleh project lain (`be-parkir_app`)
- **Port 5433** digunakan untuk UNSRI Backend PostgreSQL
- Services di Docker tetap pakai port 5432 karena connect via Docker network
- Hanya migrations dan connection dari host yang pakai port 5433

