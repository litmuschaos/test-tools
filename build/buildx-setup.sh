#!/bin/bash
set -e
echo '{"experimental":true}' | sudo tee /etc/docker/daemon.json
export DOCKER_CLI_EXPERIMENTAL=enabled
docker run --rm --privileged docker/binfmt:a7996909642ee92942dcd6cff44b9b95f08dad64
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update && sudo apt-get -y -o Dpkg::Options::="--force-confnew" install docker-ce=5:19.03.8~3-0~ubuntu-bionic # pin version for reproducibility
docker buildx create --name multibuilder --use

