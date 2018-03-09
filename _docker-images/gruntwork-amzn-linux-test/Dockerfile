# TODO: Is it worth referencing a specific tag instead of latest?
FROM amazonlinux:latest

# Reduce Docker image size per https://blog.replicated.com/refactoring-a-dockerfile-for-image-size/
# - The last line upgrades pip to the latest version.
RUN yum update -y && \
    yum upgrade -y && \
    yum install -y \
        jq \
        rsyslog \
        sudo \
        vim \
        wget && \
        yum clean all