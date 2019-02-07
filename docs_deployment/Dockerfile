FROM alpine:3.8

ENV CHAMBER_VERSION=2.2.0

COPY start_website_build /usr/bin/start_website_build

RUN apk update && apk upgrade

RUN apk --no-cache add bash git openssh curl jq && \
    curl -OL https://github.com/segmentio/chamber/releases/download/v${CHAMBER_VERSION}/chamber-v${CHAMBER_VERSION}-linux-amd64 && \
    mv chamber-v${CHAMBER_VERSION}-linux-amd64 /usr/bin/chamber && \
    chmod 755 /usr/bin/chamber && \
    chmod 755 /usr/bin/start_website_build

CMD ["/bin/bash"]