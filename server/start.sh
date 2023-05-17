#!/bin/sh

if [ "${ENV}" = "development" ]; then
    cd cmd && exec air
else
    exec ../main/main
fi
