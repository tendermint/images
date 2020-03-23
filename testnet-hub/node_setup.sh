#!/bin/bash

source /tmp/files/set_env.sh


mkdir -p "${GOPATH}"/src/github.com/cosmos
mkdir -p "${HOME}"/go/bin

cd "${GOPATH}"/src/github.com/cosmos || exit 1
git clone https://github.com/cosmos/gaia
cd gaia && git checkout "${GAIA_VERSION}"
make install

sudo cp "${GOBIN}"/gaiacli /usr/bin/gaiacli
sudo cp "${GOBIN}"/gaiad /usr/bin/gaiad

sudo rm -rf "${GOPATH}"

sudo mv /tmp/files/gaiacli.service /lib/systemd/system/gaiacli.service
sudo mv /tmp/files/gaiad.service /lib/systemd/system/gaiad.service
sudo systemctl daemon-reload

sudo gaiad init ubuntu-node --home /home/gaia/.gaiad
sudo mv /tmp/files/genesis.json /home/gaia/.gaiad/config/genesis.json
sudo chown -R gaia:gaia /home/gaia/.gaiad

#### NGINX setup
####
sudo mv /tmp/files/nginx/nginx.conf /etc/nginx/nginx.conf
sudo mv /tmp/files/nginx/sites-available/gaia* /etc/nginx/sites-available/

sudo sed -i -- 's/<DNS-NAME>/'"$GAIAD_DNS_NAME"'/g' /etc/nginx/sites-available/gaiad
sudo sed -i -- 's/<DNS-NAME>/'"$GAIACLI_DNS_NAME"'/g' /etc/nginx/sites-available/gaiacli

sudo ln -s /etc/nginx/sites-available/gaiacli /etc/nginx/sites-enabled/gaiacli
sudo ln -s /etc/nginx/sites-available/gaiad /etc/nginx/sites-enabled/gaiad
