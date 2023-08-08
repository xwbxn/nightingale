#!/bin/bash
if [ ! -d "./public" ]; then
    mkdir -p public/prod-api
    cp -r ~/ops/ythyw/build/* ./public
fi


# Embed files into a go binary
if ! statik -src=./public -dest=./front -ns=dashboard -p statik_dashboard; then
    echo "failed to embed files into a go binary!"
    exit 4
fi

rm -rf ./public