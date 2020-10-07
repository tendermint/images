FROM golang:1.15.0-buster
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get --no-install-recommends -y install \
    pciutils build-essential git wget \
    lsb-release dpkg-dev curl bsdmainutils fakeroot
RUN mkdir -p /usr/local/share/cosmos-sdk/

# Deploy the shell functions library.
COPY buildlib.sh /usr/local/share/cosmos-sdk/

# Create the 'builder' user.
RUN useradd -ms /bin/bash -U builder
ARG APP
ARG DEBUG
ENV APP ${APP:-cosmos-sdk}
ENV DEBUG ${DEBUG}
ENV VERSION unknown
ENV COMMIT unknown
ENV LEDGER_ENABLE true

# Drop root privileges.
USER builder:builder
WORKDIR /sources

# Mount the application's sources.
VOLUME [ "/sources" ]

# Run the application's build.sh.
ENTRYPOINT [ "/sources/build.sh" ]
