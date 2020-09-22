FROM golang:1.14-alpine3.11 as builder

RUN apk add --update --no-cache build-base curl git upx && \
  rm -rf /var/cache/apk/*

ENV GOLANG_PROTOBUF_VERSION=1.3.1 \
  GOGO_PROTOBUF_VERSION=1.3.1
RUN GO111MODULE=on go get \
  github.com/golang/protobuf/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} \
  github.com/gogo/protobuf/protoc-gen-gogo@v${GOGO_PROTOBUF_VERSION} \
  github.com/gogo/protobuf/protoc-gen-gogofaster@v${GOGO_PROTOBUF_VERSION} && \
  mv /go/bin/protoc-gen-go* /usr/local/bin/

ENV PROTOBUF_VERSION=3.11.4
RUN mkdir -p /tmp/protoc && \
  curl -sSL \
  https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip \
  -o /tmp/protoc/protoc.zip && \
  cd /tmp/protoc && \
  unzip protoc.zip && \
  mv /tmp/protoc/include /usr/local/include

RUN upx --lzma /usr/local/bin/*

FROM alpine:edge

WORKDIR /work

ENV \
  GRPC_VERSION=1.32.0 \
  PROTOBUF_VERSION=3.13.0 \
  ALPINE_GRPC_VERSION_SUFFIX=r0 \
  ALPINE_PROTOBUF_VERSION_SUFFIX=r0

RUN echo 'http://dl-cdn.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
  apk add --update --no-cache grpc=${GRPC_VERSION}-${ALPINE_GRPC_VERSION_SUFFIX} protobuf=${PROTOBUF_VERSION}-${ALPINE_PROTOBUF_VERSION_SUFFIX} clang && \
  rm -rf /var/cache/apk/*

COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/local/include /usr/include

ENV GOGO_PROTOBUF_VERSION=1.2.1 \
  GOLANG_PROTOBUF_VERSION=1.3.1 \
  GRPC_GATEWAY_VERSION=1.8.5 \
  GRPC_WEB_VERSION=1.0.4 \
  TWIRP_VERSION=5.7.0 \
  YARPC_VERSION=1.37.3
