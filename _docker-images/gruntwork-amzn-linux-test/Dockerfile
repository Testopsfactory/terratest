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

# Installing pip with yum doesn't actually put it in the PATH, so we use easy_install instead. Pip will now be placed
# in /usr/local/bin, but amazonlinux's sudo uses a sanitzed PATH that does not include /usr/local/bin, so we symlink pip.
RUN curl https://bootstrap.pypa.io/ez_setup.py | sudo /usr/bin/python && \
    easy_install pip && \
    pip install --upgrade pip && \
    ln -s /usr/local/bin/pip /usr/bin/pip

# Install the AWSCLI (which apparently does not come pre-bundled with Amazon Linux!)
RUN pip install --upgrade setuptools && \
    pip install awscli --upgrade
