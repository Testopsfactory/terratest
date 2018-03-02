FROM ubuntu:16.04

# Reduce Docker image size per https://blog.replicated.com/refactoring-a-dockerfile-for-image-size/
# - The last line upgrades pip to the latest version.
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    jq \
    python-pip \
    rsyslog \
    software-properties-common \
    sudo \
    wget && \
    rm -rf /var/lib/apt/lists/* && \
    pip install --upgrade setuptools && \
    pip install --upgrade pip

# Install the AWS CLI per https://docs.aws.amazon.com/cli/latest/userguide/installing.html
RUN pip install awscli --upgrade