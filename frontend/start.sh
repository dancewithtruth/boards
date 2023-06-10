#!/bin/sh

if [ "${ENV}" = "development" ]; then
    npm run dev
else
    npm run build && npm run start
fi
