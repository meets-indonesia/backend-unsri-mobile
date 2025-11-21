#!/bin/bash

# Script to run all services in separate terminal windows/tabs
# Usage: ./scripts/run-all-services.sh

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting all UNSRI Backend Services...${NC}"
echo ""

# Check if database and Redis are running
echo -e "${BLUE}Checking prerequisites...${NC}"
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "⚠️  PostgreSQL is not running. Please start it first:"
    echo "   docker-compose -f deployments/docker-compose/docker-compose.yml up -d postgres"
    exit 1
fi

if ! redis-cli ping > /dev/null 2>&1; then
    echo "⚠️  Redis is not running. Please start it first:"
    echo "   docker-compose -f deployments/docker-compose/docker-compose.yml up -d redis"
    exit 1
fi

echo -e "${GREEN}✓ PostgreSQL and Redis are running${NC}"
echo ""

# Function to run service in new terminal (macOS)
run_in_terminal() {
    local service_name=$1
    local command=$2
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        osascript -e "tell application \"Terminal\" to do script \"cd $(pwd) && echo 'Starting $service_name...' && $command\""
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        gnome-terminal -- bash -c "cd $(pwd) && echo 'Starting $service_name...' && $command; exec bash"
    else
        echo "Unsupported OS. Please run services manually in separate terminals."
        exit 1
    fi
}

# Run each service in separate terminal
echo -e "${BLUE}Starting services in separate terminals...${NC}"

run_in_terminal "API Gateway" "go run ./cmd/api-gateway"
sleep 1

run_in_terminal "Auth Service" "go run ./cmd/auth-service"
sleep 1

run_in_terminal "User Service" "go run ./cmd/user-service"
sleep 1

run_in_terminal "Attendance Service" "go run ./cmd/attendance-service"
sleep 1

run_in_terminal "Schedule Service" "go run ./cmd/schedule-service"
sleep 1

run_in_terminal "QR Service" "go run ./cmd/qr-service"
sleep 1

run_in_terminal "Course Service" "go run ./cmd/course-service"
sleep 1

run_in_terminal "Broadcast Service" "go run ./cmd/broadcast-service"
sleep 1

run_in_terminal "Notification Service" "go run ./cmd/notification-service"
sleep 1

run_in_terminal "Calendar Service" "go run ./cmd/calendar-service"

echo ""
echo -e "${GREEN}All services started!${NC}"
echo ""
echo "Services are running on:"
echo "  - API Gateway: http://localhost:8080"
echo "  - Auth Service: http://localhost:8081"
echo "  - User Service: http://localhost:8082"
echo "  - Schedule Service: http://localhost:8083"
echo "  - Attendance Service: http://localhost:8084"
echo "  - QR Service: http://localhost:8085"
echo "  - Broadcast Service: http://localhost:8086"
echo "  - Notification Service: http://localhost:8087"
echo "  - Calendar Service: http://localhost:8088"
echo "  - Course Service: http://localhost:8089"
echo ""
echo "To stop services, close the terminal windows or press Ctrl+C in each terminal."

