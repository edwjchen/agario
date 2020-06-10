#!/bin/bash
cd 
sudo rm /var/lib/dpkg/lock-frontend
sudo rm /var/lib/dpkg/lock

curl -O https://bootstrap.pypa.io/get-pip.py
sudo apt install python3-distutils -y
python3 get-pip.py
echo "export PATH=/home/ubuntu/.local/bin:$PATH" >> ~/.bashrc
echo "export SDL_VIDEODRIVER=dummy" >> ~/.bashrc
source ~/.bashrc
pip3 install grpcio pygame protobuf

echo "export GOPATH=$HOME/agario" >> ~/.bashrc
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc

cd agario/src
cd peer_to_peer
go get -v ./...
cd ../client_server/server
go get -v ./...

cd

