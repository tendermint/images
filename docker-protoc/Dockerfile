FROM ubuntu:18.04

ENV PROTOC_VERSION=3.11.2
ENV GO_VERSION=1.13.5
ENV BUF_VERSION=0.5.0
ENV PATH="${PATH}:/usr/local/go/bin:${HOME}/go/bin"

RUN apt-get update && apt-get install -y --no-install-recommends curl git make unzip ca-certificates
RUN rm -rf /var/lib/apt/lists/*

## Install protoc
RUN curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip"

RUN unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local bin/protoc
RUN unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local 'include/*'
RUN rm -rf protoc-${PROTOC_VERSION}-linux-x86_64.zip

## Install golang

RUN curl -OL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"

## Install buf

RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" -o "${BIN}/buf"
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-breaking-Linux-x86_64" -o "${BIN}/protoc-gen-buf-check-breaking"
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-lint-Linux-x86_64" -o "${BIN}/protoc-gen-buf-check-lint"
RUN chmod +x "${BIN}/buf" "${BIN}/protoc-gen-buf-check-breaking" "${BIN}/protoc-gen-buf-check-lint"

CMD ["/bin/bash"]
