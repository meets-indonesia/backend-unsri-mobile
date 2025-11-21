#!/bin/bash

# Database migration script
# Usage: ./scripts/migrate.sh [up|down]

MIGRATION_DIR="migrations"
DATABASE_URL="postgres://unsri_user:unsri_pass@localhost:5432/unsri_db?sslmode=disable"

if [ -z "$1" ]; then
    echo "Usage: $0 [up|down]"
    exit 1
fi

if [ "$1" = "up" ]; then
    echo "Running migrations up..."
    migrate -path $MIGRATION_DIR -database "$DATABASE_URL" up
elif [ "$1" = "down" ]; then
    echo "Running migrations down..."
    migrate -path $MIGRATION_DIR -database "$DATABASE_URL" down
else
    echo "Invalid argument: $1"
    echo "Usage: $0 [up|down]"
    exit 1
fi

