.PHONY: build run test clean migrate-up migrate-down docker-build docker-up docker-down k8s-deploy k8s-delete

# Build all services
build:
	@echo "Building all services..."
	@go build -o bin/api-gateway ./cmd/api-gateway
	@go build -o bin/auth-service ./cmd/auth-service
	@go build -o bin/user-service ./cmd/user-service
	@go build -o bin/attendance-service ./cmd/attendance-service
	@go build -o bin/course-service ./cmd/course-service
	@go build -o bin/broadcast-service ./cmd/broadcast-service
	@go build -o bin/notification-service ./cmd/notification-service
	@go build -o bin/schedule-service ./cmd/schedule-service
	@go build -o bin/qr-service ./cmd/qr-service
	@go build -o bin/calendar-service ./cmd/calendar-service
	@go build -o bin/location-service ./cmd/location-service
	@go build -o bin/access-service ./cmd/access-service
	@go build -o bin/quick-actions-service ./cmd/quick-actions-service
	@go build -o bin/file-storage-service ./cmd/file-storage-service
	@go build -o bin/search-service ./cmd/search-service
	@go build -o bin/report-service ./cmd/report-service
	@echo "Build complete!"

run-user-service:
	@go run ./cmd/user-service

# Run services locally
run-api-gateway:
	@go run ./cmd/api-gateway

run-auth-service:
	@go run ./cmd/auth-service

run-user-service:
	@go run ./cmd/user-service

run-attendance-service:
	@go run ./cmd/attendance-service

run-course-service:
	@go run ./cmd/course-service

run-broadcast-service:
	@go run ./cmd/broadcast-service

run-notification-service:
	@go run ./cmd/notification-service

run-schedule-service:
	@go run ./cmd/schedule-service

run-qr-service:
	@go run ./cmd/qr-service

run-calendar-service:
	@go run ./cmd/calendar-service

run-qr-service:
	@go run ./cmd/qr-service

run-location-service:
	@go run ./cmd/location-service

run-access-service:
	@go run ./cmd/access-service

run-quick-actions-service:
	@go run ./cmd/quick-actions-service

run-file-storage-service:
	@go run ./cmd/file-storage-service

run-search-service:
	@go run ./cmd/search-service

run-report-service:
	@go run ./cmd/report-service

swagger:
	@swag init -g cmd/api-gateway/main.go
	@echo "Swagger docs generated in docs/ folder"

# Run all services in background (requires multiple terminals)
run-all:
	@echo "Starting all services..."
	@echo "Please run each service in separate terminals:"
	@echo "  Terminal 1: make run-api-gateway"
	@echo "  Terminal 2: make run-auth-service"
	@echo "  Terminal 3: make run-user-service"
	@echo "  Terminal 4: make run-attendance-service"
	@echo "  Terminal 5: make run-schedule-service"
	@echo "  Terminal 6: make run-qr-service"
	@echo "  Terminal 7: make run-course-service"
	@echo "  Terminal 8: make run-broadcast-service"
	@echo "  Terminal 9: make run-notification-service"
	@echo "  Terminal 10: make run-calendar-service"

# Run tests
test:
	@go test -v ./...

# Clean build artifacts
clean:
	@rm -rf bin/
	@echo "Clean complete!"

# Database migrations
migrate-up:
	@migrate -path migrations -database "postgres://unsri_user:unsri_pass@localhost:5432/unsri_db?sslmode=disable" up

migrate-down:
	@migrate -path migrations -database "postgres://unsri_user:unsri_pass@localhost:5432/unsri_db?sslmode=disable" down

# Docker commands
docker-build:
	@docker-compose -f deployments/docker-compose/docker-compose.yml build

docker-up:
	@docker-compose -f deployments/docker-compose/docker-compose.yml up -d

docker-down:
	@docker-compose -f deployments/docker-compose/docker-compose.yml down

docker-logs:
	@docker-compose -f deployments/docker-compose/docker-compose.yml logs -f

# Kubernetes commands
k8s-deploy:
	@kubectl apply -f deployments/kubernetes/

k8s-delete:
	@kubectl delete -f deployments/kubernetes/

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run

