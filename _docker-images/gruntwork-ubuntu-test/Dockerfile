FROM ubuntu:16.04

# Reduce Docker image size per https://blog.replicated.com/refactoring-a-dockerfile-for-image-size/
# - dnsutils: Install handy DNS checking tools like dig
# - libcrypt-hcesha-perl: Install shasum
# - software-properties-common: Install add-apt-repository
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install --no-install-recommends -y \
    ca-certificates \
    curl \
    dnsutils \
    jq \
    libcrypt-hcesha-perl \
    python-pip \
    rsyslog \
    software-properties-common \
    sudo \
    vim \
    wget && \
    rm -rf /var/lib/apt/lists/*

# Install the AWS CLI per https://docs.aws.amazon.com/cli/latest/userguide/installing.html. The last line upgrades pip
# to the latest version.
RUN pip install --upgrade setuptools && \
    pip install --upgrade pip && \
    pip install awscli --upgrade