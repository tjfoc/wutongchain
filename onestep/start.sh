#!/bin/bash
source ./stop.sh

echo "=====copying exec to Dockerfiles===="
cp -p ./peer/peer ./Dockerfiles/
cp -p ./sdk/sdk ./Dockerfiles/

chmod a+rw /var/run/docker.sock
mkdir -p ./peer/Peer1/logs/
mkdir -p ./peer/Peer2/logs/
chmod a+rwx ./peer/Peer1/logs/
chmod a+rwx ./peer/Peer2/logs/
echo "===== get ccenv image ===="
docker pull tjfoc/tjfoc-ccenv:1.0.1
echo "=====build peer exec image===="
docker build -f ./Dockerfiles/DockerFilePeer -t peer-image:latest ./Dockerfiles
echo "=====build sdk exec image===="
docker build -f ./Dockerfiles/DockerFileSdk -t sdk-image:latest ./Dockerfiles
echo "=====list images====="
docker images

echo "=====start docker-compose====="
docker-compose -f config.yaml up


