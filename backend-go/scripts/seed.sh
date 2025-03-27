#!/bin/bash

cd "$(dirname "$0")/.."

echo "Seeding database..."
ENV=dev go run ../backend-go/cmd/migrate/seed/main.go