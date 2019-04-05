#!/usr/bin/bash

export PATH=$PATH:/usr/local/go/bin

mkdir -p $HOME/go/bin

export GOPATH=$HOME/go
export GOBIN=${GOPATH}/bin
export PATH=$PATH:${GOBIN}

mkdir -p ${GOPATH}/src/github.com/cosmos
cd ${GOPATH}/src/github.com/cosmos
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout master
make tools install

cp /tmp/notify_slack.go ~/go/src/github.com/cosmos/cosmos-sdk/notify_slack.go
cp /tmp/multisim.sh ~/go/src/github.com/cosmos/cosmos-sdk/multisim.sh
