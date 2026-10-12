#!/bin/bash
set -e

cd /app

# go.sum belum ada -> start pertama: resolve dependency (fiber, gorm, mysql, uuid)
if [ ! -f "go.sum" ]; then
    echo "First run: go mod tidy ..."
    go mod tidy
else
    go mod download
fi

mkdir -p tmp

echo ""
echo "======================================"
echo " Go Fiber API  (hot reload: Air)"
echo "======================================"
echo "DB   : ${DB_USERNAME}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}"
echo "Port : ${APP_PORT:-8000} (container)"
echo ""

exec "$@"
