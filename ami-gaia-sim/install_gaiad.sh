#!/usr/bin/bash

sudo mv /tmp/set_env.sh /etc/profile.d/set_env.sh
sudo chmod u+x /etc/profile.d/set_env.sh
source /etc/profile.d/set_env.sh

mkdir -p ${HOME}/go/bin
mkdir -p ${GOPATH}/src/github.com/cosmos

cd ${GOPATH}/src/github.com/cosmos
git clone https://github.com/cosmos/cosmos-sdk
cd cosmos-sdk && git checkout ${GAIA_COMMIT_HASH}
make runsim tools

sudo mv /home/ec2-user/go/bin/runsim /usr/bin/runsim
