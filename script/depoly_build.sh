#!/usr/bin/env bash
RUN_NAME="SDAS"
go build -ldflags="-s -w" -o output/bin/${RUN_NAME} && upx -9 output/bin/${RUN_NAME}