#!/usr/bin/env bash
RUN_NAME="SDAS"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o output/bin/${RUN_NAME}_linux && build/upx -9 output/bin/${RUN_NAME}_linux
echo "exec build success"
docker build -t sdas .
echo "docker build success"
