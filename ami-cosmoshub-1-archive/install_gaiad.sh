#!/usr/bin/bash

export PATH=$PATH:/usr/local/go/bin

mkdir -p "${HOME}"/go/bin

export GOPATH=$HOME/go
export GOBIN=${GOPATH}/bin
export PATH=$PATH:${GOBIN}

mkdir -p "${GOPATH}"/src/github.com/cosmos
cd "${GOPATH}"/src/github.com/cosmos || exit 1
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout "${GAIAD_VERSION}"
make tools install

sudo mv "${GOBIN}"/gaiad /usr/bin/gaiad
sudo mv /tmp/mount_ebs.sh /usr/bin/mount_ebs.sh
sudo chmod u+x /usr/bin/mount_ebs.sh

# Configure gaiad service
sudo mv /tmp/gaiad.service /usr/lib/systemd/system/gaiad.service
sudo systemctl enable gaiad

sudo rm -rf "${GOPATH}"