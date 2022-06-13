#!/bin/bash
# This file contains all installed packages and their versions

# Functions
checkPkg() {
    if ! command -v $pkg &>/dev/null; then
        sudo apt-get install $pkg
    else
        echo "Package exist"
    fi
}

# --------------------------Install node and npm--------------------------

# Install nodejs and npm:
curl -sL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# last install 02.march.2022:  => node version: v16.14.0
# last install 02.march.2022:  => npm version: 8.3.1

# --------------------------Install Golang--------------------------
# install gcc and g++ compilers and libraries
apt-get install build-essential

# download and Install golang version 1.17.7
curl -OL https://golang.org/dl/go1.17.7.linux-amd64.tar.gz

# check the integrity
# Run:
sha256sum go1.17.7.linux-amd64.tar.gz
# If the checksum matches the one listed on the downloads page, you’ve done this step correctly.

sudo tar -C /usr/local -xvf go1.17.7.linux-amd64.tar.gz

# set the path to profile
echo 'export PATH=$PATH:/usr/local/go/bin' >>~/.bashrc
source ~/.bashrc

# -------------------------Install python pkgs--------------------------
checkPkg python3-pip

# --------------------------Enable SSH to github--------------------------
# Check if keypair exist

File1=~/.ssh/id_ed25519.pub
if [ -if "$File"]; then
    echo "$File exist."
else
    ssh-keygen -t ed25519 -C "mohamad.h@hotmail.no"
    # Enter a file in which to save the key (/home/you/.ssh/algorithm): ->[Press enter]
    # Enter passphrase (empty for no passphrase): ->[Type a passphrase]
fi
#start

# Adding your SSH key to the ssh-agent
# First: Start the ssh-agent in the background.
eval "$(ssh-agent -s)"
# Then: Add your SSH private key to the ssh-agent
ssh-add ~/.ssh/id_ed25519

# Copy the publick key and add to github account:
cat ~/.ssh/id_ed25519.pub
# setting -> under access: ssh key .. -> add new key
echo "Have you added the key to your github account? [if yes press enter]"
read answer
# verify connection
ssh -T git@github.com

# --------------------------Fabric test netowork--------------------------
# prerequests
# Install git if not installed:
checkPkg git

# Install curl if not installed:
checkPkg curl

# Install the latest version
# sudo apt-get -y install docker-compose
checkPkg docker-compose

# verify docker and docker compose
docker --version
docker-compose --version

# run docker deamon
sudo systemctl start docker

# start docker by os sturtup
sudo systemctl enable docker

# check if running: sudo systemctl status docker

# add your user to docker
sudo groupadd docker
sudo gpasswd -a $USER docker

# Install jq
checkPkg jq
chmod +x jq

# download the Fabric samples, docker images, and binaries
mkdir -p $HOME/thesis
cd $HOME/thesis

# OR: if you want to put the test network and binaries in another repo/dir go directly to the next step

# the latest release of Fabric samples, docker images, and binaries
curl -sSL https://bit.ly/2ysbOFE | bash -s

#--------------------------.NET------------------------- ** here

# make a simple demo
cd $HOME/DID/aries-test/AriesAgent
dotnet add package AgentFramework --version 4.0.1
dotnet add package Hyperledger.Aries.AspNetCore --version 1.6.2

# Install dotnet
wget https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
sudo dpkg -i packages-microsoft-prod.deb
rm packages-microsoft-prod.deb

# install dotnet-sdk
sudo apt-get update
sudo apt-get install -y apt-transport-https &&
    sudo apt-get update &&
    sudo apt-get install -y dotnet-sdk-6.0

#--------------------------Install Indy-SDK------------------------- *** WORKS ONLY IN UBUNTU 16.X OR 18.X, NOT IN 20.X version(s)
# First OPTION:
# THIS CAN BE DONE IN ANOTHER HOST (ubuntu 18.x) IN THE SAME NETWORK.
# In our case the IP address of other host is: 192.168.0.28

# Install Indy-SDK *** NOT WORKING CURRENTLY*************
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys CE7709D068DB5E88
sudo add-apt-repository "deb https://repo.sovrin.org/sdk/deb bionic stable"
sudo apt-get update
sudo apt-get install -y {library}

#--------------------------OR= VON network (Indy)
# https://github.com/bcgov/von-network/blob/main/docs/UsingVONNetwork.md

git clone https://github.com/bcgov/von-network
cd von-network
./manage build
./manage start --logs

#--------------------------Install Rust--------------------------

# install rust and rustup
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
# verify
rustc --version

#--------------------------START ACA-Py in docker-OPEN API sample--------------------------

#install ngrok
sudo apt update
checkPkg snapd
sudo snap install ngrok

# DO the following steps for each agents
# navgigate to your dir
cd ~/your_dir
git clone https://github.com/hyperledger/aries-cloudagent-python
cd aries-cloudagent-python/demo
# enable logging
mkdir ../logs
chmod uga+rws ../logs
# IMPORTANT ADD THE FLAG : --network=host TO the line 150 in run_demo (for building the Dockerfile.demo)
LEDGER_URL=http://dev.greenlight.bcovrin.vonx.io ./run_demo faber --events --no-auto --bg
LEDGER_URL=http://dev.greenlight.bcovrin.vonx.io ./run_demo alice --events --no-auto --bg

#--------------------------Install protobuf compiler--------------------------
sudo apt update
sudo apt install protobuf-compiler

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

export PATH="$PATH:$(go env GOPATH)/bin"
