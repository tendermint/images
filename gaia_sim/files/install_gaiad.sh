#!/usr/bin/bash

sudo mv /tmp/set_env.sh /etc/profile.d/set_env.sh
sudo chmod u+x /etc/profile.d/set_env.sh
source /etc/profile.d/set_env.sh

mkdir -p ${HOME}/go/bin
mkdir -p ${GOPATH}/src/github.com/cosmos

cd ${GOPATH}/src/github.com/cosmos
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout master
make tools install

mv /tmp/notify_slack.go ~/go/src/github.com/cosmos/cosmos-sdk/notify_slack.go
mv /tmp/multisim.sh ~/go/src/github.com/cosmos/cosmos-sdk/multisim.sh
chmod u+x ~/go/src/github.com/cosmos/cosmos-sdk/multisim.sh
