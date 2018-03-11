# Gruntwork CentOS-Test Docker Image

The purpose of this Docker image is to provide a pre-built CentOS 7 Docker image that has most of the libraries
we would expect to be installed on the CentOS 7 AMI that would run in AWS. For example, we'd expect `sudo` in AWS,
but it doesn't exist by default in Docker `centos:7`.

### Building and Pushing a New Docker Image to Docker Hub

This Docker image should publicly accessible via Docker Hub at https://hub.docker.com/r/gruntwork/centos-test/. To build and
upload it:

1. `docker build -t gruntwork/centos-test:7 .`
1. `docker push gruntwork/centos-test:7`

