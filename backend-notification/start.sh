#!/bin/sh

if [ "$ENV" = "development" ]; then
    echo "Starting notification service using hot reload"
    air
else
    echo "Starting notification service using build binary"
    /app/main
fi
