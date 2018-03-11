FROM centos:7

# Reduce Docker image size per https://blog.replicated.com/refactoring-a-dockerfile-for-image-size/
# - The last line upgrades pip to the latest version.
RUN yum update -y && \
    yum upgrade -y && \
    yum install -y epel-release && \
    yum install -y \
        bind-utils \
        python-pip \
        rsyslog \
        sudo \
        vim \
        wget && \
        yum clean all

# Install jq. Oddly, there's no RPM for jq, so we install the binary directly. https://serverfault.com/a/768061/199943
RUN wget -O jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
    chmod +x ./jq && \
    cp jq /usr/bin

# Install the AWS CLI per https://docs.aws.amazon.com/cli/latest/userguide/installing.html
RUN pip install --upgrade pip && \
    pip install --upgrade setuptools && \
    pip install awscli --upgrade