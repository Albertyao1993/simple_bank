#!/bin/sh

set -e

echo "run db migration"
source /app/app.env
/usr/bin/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec /app/main