#!/usr/bin/env bash
RUN_NAME="SDAS"

mkdir -p output/bin
mkdir -p output/image
cp script/* output/
chmod +x output/bootstrap.sh

GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o output/bin/${RUN_NAME}_linux && build/upx -9 output/bin/${RUN_NAME}_linux
echo "exec build success"
docker build -t sdas:v0.0.6 .
echo "docker build success"
docker save -o output/image/SDAS.tar sdas:latest
cp script/docker-install.sh output/image
zip -r output/image/SDAS.zip output/image/
rm output/image/docker-install.sh