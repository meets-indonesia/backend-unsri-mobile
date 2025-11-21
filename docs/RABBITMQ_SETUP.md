# RabbitMQ Setup & Usage

## Overview

RabbitMQ digunakan untuk:
1. **Scheduled Broadcasts** - Publish broadcast yang dijadwalkan
2. **Async Notifications** - Send notifications secara asynchronous
3. **Event-driven Communication** - Komunikasi antar services

## Setup

### 1. Docker Compose

RabbitMQ sudah ditambahkan ke `docker-compose.yml`:

```yaml
rabbitmq:
  image: rabbitmq:3-management-alpine
  container_name: unsri-rabbitmq
  ports:
    - "5672:5672"   # AMQP port
    - "15672:15672" # Management UI
  environment:
    RABBITMQ_DEFAULT_USER: unsri_user
    RABBITMQ_DEFAULT_PASS: unsri_pass
```

### 2. Start RabbitMQ

```bash
cd deployments/docker-compose
docker-compose up -d rabbitmq
```

### 3. Access Management UI

- URL: http://localhost:15672
- Username: `unsri_user`
- Password: `unsri_pass`

## Usage

### Shared Package

Gunakan `internal/shared/messaging/rabbitmq.go` untuk koneksi RabbitMQ:

```go
import "unsri-backend/internal/shared/messaging"

// Initialize
rabbitmq, err := messaging.NewRabbitMQ(messaging.Config{
    Host:     "localhost",
    Port:     "5672",
    User:     "unsri_user",
    Password: "unsri_pass",
    VHost:    "/",
})
```

### Queues

#### 1. Scheduled Broadcasts Queue
- **Queue Name**: `scheduled_broadcasts`
- **Exchange**: `broadcasts`
- **Routing Key**: `scheduled`
- **Purpose**: Publish broadcast yang dijadwalkan

#### 2. Notifications Queue
- **Queue Name**: `notifications`
- **Exchange**: `notifications`
- **Routing Key**: `send`
- **Purpose**: Send notifications secara async

## Implementation Examples

### Broadcast Service - Scheduled Broadcasts

```go
// Publish scheduled broadcast
rabbitmq.Publish(
    "broadcasts",
    "scheduled",
    false,
    false,
    amqp.Publishing{
        ContentType: "application/json",
        Body:        []byte(broadcastJSON),
    },
)
```

### Notification Service - Async Notifications

```go
// Publish notification
rabbitmq.Publish(
    "notifications",
    "send",
    false,
    false,
    amqp.Publishing{
        ContentType: "application/json",
        Body:        []byte(notificationJSON),
    },
)
```

## Environment Variables

Tambahkan ke service config:

```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=unsri_user
RABBITMQ_PASSWORD=unsri_pass
RABBITMQ_VHOST=/
```

## Monitoring

- Management UI: http://localhost:15672
- Check queues, exchanges, connections
- Monitor message rates

