#!/bin/sh

migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

if [ "${ENV}" = "development" ]; then
    exec air
else
    exec ./cmd/main
fi
