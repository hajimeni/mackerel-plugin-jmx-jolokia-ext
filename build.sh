#!/usr/bin/env bash

OS=("darwin" "linux" "windows")
ARCH="amd64"
NAME=$(basename $(pwd))

for o in ${OS[@]}; do
    GOOS=$o GOARCH=$ARCH go build -v -a -o bin/$NAME-$o-$ARCH
done
