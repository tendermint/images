FROM alpine:3.8

ENV TERRAFORM_VERSION=0.11.11
ENV TERRAGRUNT_VERSION=0.17.3
ENV CHAMBER_VERSION=2.2.0

COPY config /root/.aws/
COPY known_hosts /root/.ssh/known_hosts

RUN apk update && apk upgrade

RUN apk --no-cache add bash curl openssh git unzip jq python3 outils-cksum && \
    python3 -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip3 install --upgrade pip setuptools && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
    if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
    pip install awscli --upgrade && \
    curl -OL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    curl -OL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS && \
    cat terraform_${TERRAFORM_VERSION}_SHA256SUMS | grep linux_amd64 | sha256sum -c -  && \
    curl -OL https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64 && \
    curl -OL https://github.com/segmentio/chamber/releases/download/v${CHAMBER_VERSION}/chamber-v${CHAMBER_VERSION}-linux-amd64 && \
    curl -OL https://github.com/segmentio/chamber/releases/download/v${CHAMBER_VERSION}/chamber-v${CHAMBER_VERSION}.sha256sums && \
    cat chamber-v${CHAMBER_VERSION}.sha256sums | grep linux-amd64 | sha256sum -c - && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin/ && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terragrunt_linux_amd64 /usr/bin/terragrunt && \
    mv chamber-v${CHAMBER_VERSION}-linux-amd64 /usr/bin/chamber && \
    chmod 755 /usr/bin/terraform && \
    chmod 755 /usr/bin/terragrunt && \
    chmod 755 /usr/bin/chamber && \
    chmod 644 /root/.ssh/known_hosts && \
    rm -r /root/.cache

CMD ["/bin/bash"]