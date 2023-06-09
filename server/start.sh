#!/bin/sh

if [ "${ENV}" = "development" ]; then
    exec air
else
    exec ./cmd/main
fi
