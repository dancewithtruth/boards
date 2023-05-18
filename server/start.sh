#!/bin/sh

migrate -path ./db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@db:${DB_PORT}/${DB_NAME}?sslmode=disable" up

if [ "${ENV}" = "development" ]; then
    exec air
else
    exec ../main/main
fi
