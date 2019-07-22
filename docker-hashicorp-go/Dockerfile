FROM golang:1.12.7-stretch

ENV TERRAFORM_VERSION=0.12.5
ENV PACKER_VERSION=1.4.2
ENV TERRAGRUNT_VERSION=0.19.9
ENV GO111MODULE=on

RUN apt-get -y update && \
    apt-get -y install bash zip unzip git curl jq make && \
    curl -OL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    curl -OL https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip && \
    curl -OL https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64 && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin/ && \
    unzip packer_${PACKER_VERSION}_linux_amd64.zip -d /usr/bin && \
    mv terragrunt_linux_amd64 /usr/bin/terragrunt && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    rm packer_${PACKER_VERSION}_linux_amd64.zip && \
    chmod 755 /usr/bin/terraform && \
    chmod 755 /usr/bin/packer && \
    chmod 755 /usr/bin/terragrunt

RUN go get github.com/aws/aws-sdk-go

CMD ["/bin/bash"]