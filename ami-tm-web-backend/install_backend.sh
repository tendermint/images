#!/bin/bash -xe

git clone https://github.com/tendermint/backend.git

cd backend

npm install

cd ..

sudo mv backend /home/backend
sudo chown -R backend:backend /home/backend/backend

sudo mv /tmp/backend.service /lib/systemd/system/backend.service
sudo systemctl enable backend