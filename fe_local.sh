#!/bin/bash
if [ ! -d "./pub" ]; then
    cp -r ~/ops/fe/pub .
fi


# Embed files into a go binary
if ! $GOPATH/bin/statik -src=./pub -dest=./front; then
    echo "failed to embed files into a go binary!"
    exit 4
fi

rm -rf ./pub