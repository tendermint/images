#!/usr/bin/bash

export PATH=$PATH:/usr/local/go/bin

mkdir -p $HOME/go/bin

export GOPATH=$HOME/go
export GOBIN=${GOPATH}/bin
export PATH=$PATH:${GOBIN}

mkdir -p ${GOPATH}/src/github.com/cosmos
cd ${GOPATH}/src/github.com/cosmos
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout ${GAIAD_VERSION}
make tools install

sudo cp ${GOBIN}/gaiacli /usr/bin/gaiacli

sudo rm -rf ${GOPATH}

sudo mv /tmp/nginx.conf /etc/nginx/nginx.conf
sudo mv /tmp/gaiacli.service /usr/lib/systemd/system/gaiacli.service
sudo systemctl enable gaiacli

# SELINUX setting that allows nginx to connect to gaiacli locally
sudo setsebool -P httpd_can_network_connect 1
sudo systemctl enable nginx