FROM alpine:3.8

ENV TERRAFORM_VERSION=0.11.11
ENV TERRAGRUNT_VERSION=0.18.0

RUN apk --no-cache add bash unzip git openssh curl jq && \
    curl -OL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    curl -OL https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64 && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin/ && \
    mv terragrunt_linux_amd64 /usr/bin/terragrunt && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    chmod 755 /usr/bin/terraform && \
    chmod 755 /usr/bin/terragrunt && \
    touch /root/.bashrc

CMD ["/bin/bash"]