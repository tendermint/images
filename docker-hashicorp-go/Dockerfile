FROM golang:1.12.4-alpine3.9

ENV TERRAFORM_VERSION=0.11.13
ENV PACKER_VERSION=1.4.0
ENV TERRAGRUNT_VERSION=0.18.3
ENV GO111MODULE=on

RUN apk --no-cache add bash unzip git openssh curl jq make && \
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