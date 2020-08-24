FROM golang:alpine

COPY config /root/.aws/

RUN apk update && apk upgrade
RUN apk --no-cache add make bash curl openssh git unzip jq python3 outils-cksum npm && \
    python3 -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip3 install --upgrade pip setuptools && \
    if [ ! -e /usr/bin/pip ]; then ln -s pip3 /usr/bin/pip ; fi && \
    if [[ ! -e /usr/bin/python ]]; then ln -sf /usr/bin/python3 /usr/bin/python; fi && \
    pip install awscli --upgrade && \
    rm -r /root/.cache

CMD ["/bin/bash"]
