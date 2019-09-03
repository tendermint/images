#!/usr/bin/bash

sudo mv /tmp/set_env.sh /etc/profile.d/set_env.sh
sudo mv /tmp/genesis.json /home/ec2-user/genesis.json
source /etc/profile.d/set_env.sh

mkdir -p "${HOME}"/go/bin
mkdir -p "${GOPATH}"/src/github.com/cosmos

cd "${GOPATH}"/src/github.com/cosmos || exit 1
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout "${GAIA_COMMIT_HASH}"
make build

cd "${GOPATH}"/src/github.com/cosmos || exit 1
git clone https://github.com/cosmos/tools
cd tools/cmd/runsim || exit 1
git checkout "${RUNSIM_COMMIT_HASH}"
go install

sudo mv /home/ec2-user/go/bin/runsim /usr/bin/runsim

