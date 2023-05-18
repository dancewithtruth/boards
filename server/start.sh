#!/bin/sh

if [ "${ENV}" = "development" ]; then
    exec air
else
    exec ../main/main
fi
