#!/bin/bash

# プロジェクトのルートに移動
cd "$(dirname "$0")/.."

echo "Running database migrations..."
ENV=dev go run backend-go/cmd/migrate/main.go
