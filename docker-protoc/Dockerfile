FROM bufbuild/buf

ENV GO_VERSION=1.14.1
ENV PATH="${PATH}:/usr/local/go/bin:/root/go/bin"
ENV BIN="/usr/local/bin"

RUN apk update && apk add curl git make unzip ca-certificates
RUN rm -rf /var/lib/apt/lists/*

## Install golang

RUN curl -OL "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"

CMD ["/bin/bash"]
