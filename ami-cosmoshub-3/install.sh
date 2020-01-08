#!/bin/bash -xe

export PATH=$PATH:/usr/local/go/bin

mkdir -p "$HOME"/go/bin

export GOPATH=$HOME/go
export GOBIN=${GOPATH}/bin
export PATH=$PATH:${GOBIN}

mkdir -p "${GOPATH}"/src/github.com/cosmos
cd "${GOPATH}"/src/github.com/cosmos
git clone https://github.com/cosmos/gaia
cd gaia && git checkout "${GAIAD_VERSION}"
make install

sudo mv "${GOBIN}"/gaiad /usr/bin/gaiad
sudo mv "${GOBIN}"/gaiacli /usr/bin/gaiacli

# Configure gaiad service
sudo mv /tmp/gaiad.service /lib/systemd/system/gaiad.service
sudo mv /tmp/gaiacli.service /lib/systemd/system/gaiacli.service
sudo mv /tmp/nginx.conf /etc/nginx/nginx.conf

sudo rm -rf "${GOPATH}"