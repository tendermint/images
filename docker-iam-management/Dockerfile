FROM alpine:3.8

ENV TERRAFORM_VERSION=0.11.13
ENV TERRAGRUNT_VERSION=0.18.3

COPY terragrunt_plan /usr/bin/terragrunt_plan

RUN apk --no-cache add bash unzip git openssh curl jq && \
    curl -OL https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    curl -OL https://github.com/gruntwork-io/terragrunt/releases/download/v${TERRAGRUNT_VERSION}/terragrunt_linux_amd64 && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin/ && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terragrunt_linux_amd64 /usr/bin/terragrunt && \
    chmod 755 /usr/bin/terraform && \
    chmod 755 /usr/bin/terragrunt && \
    chmod 755 /usr/bin/terragrunt_plan && \
    touch /root/.bashrc

CMD ["/bin/bash"]