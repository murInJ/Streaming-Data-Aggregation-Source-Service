FROM golang:1.22-alpine AS build

COPY build/ /go/src/build/
COPY client/ /go/src/client/
COPY config/ /go/src/config/
COPY kitex_gen/ /go/src/kitex_gen/
COPY script/ /go/src/script/
COPY services/ /go/src/services/
COPY utils/ /go/src/utils/
COPY go.mod go.sum *.go /go/src/

WORKDIR "/go/src/"
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=1
RUN apk update
RUN apk add pkgconf
RUN apk add --no-cache build-base
#RUN apk add --no-cache \
#		libavcodec \
#		libavutil \
#		libswscale
RUN apk add ffmpeg-dev

ENV PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig
RUN pkg-config --cflags -- libavcodec libavutil libswscale
RUN go mod tidy
RUN sh script/deploy_build.sh


FROM alpine:latest

RUN apk add ffmpeg
RUN mkdir "/app"
RUN mkdir "/app/plugins"
COPY --from=build /go/src/output/bin/SDAS_linux /app/SDAS_linux
COPY config.json /app/config.json

RUN chmod +x /app/SDAS_linux

EXPOSE 8088
WORKDIR "/app"
ENTRYPOINT ["/app/SDAS_linux"]