# Reproducible Build System

This image is meant to provide a minimal deterministic 
buildsystem for Cosmos SDK applications.

# Requirements And Usage

The client application's repository must include an
`build.sh` executable file in the root folder meant to drive the build
process. The following environment variables are passed through
and made available to the `build.sh` script:
* `APP` - the application's name.
* `VERSION` - the application's version.
* `COMMIT` - the application's VCS commit's hash.
* `TARGET_OS` - whitespace-separated list of target operating systems (`linux`, `darwin`, and `windows` are supported).
* `LEDGER_ENABLED` - whether Ledger is enabled (default: `true`).
* `DEBUG` - run build with debug output. Default: empty (disabled).

The build's outputs are produced in the top-level `artifacts` directory. An example of `build.sh` follows:

```bash
#!/bin/bash

set -ue

# Expect the following envvars to be set:
# - APP
# - VERSION
# - COMMIT
# - TARGET_OS
# - LEDGER_ENABLED
# - DEBUG

# Source builder's functions library
. /usr/local/share/cosmos-sdk/buildlib.sh

# These variables are now available
# - BASEDIR
# - OUTDIR

# Build for each os-architecture pair
for os in ${TARGET_OS} ; do
    archs="`f_build_archs ${os}`"
    exe_file_extension="`f_binary_file_ext ${os}`"
    for arch in ${archs} ; do
        make clean
        GOOS="${os}" GOARCH="${arch}" GOROOT_FINAL="$(go env GOROOT)" \
        make build \
            LDFLAGS=-buildid=${VERSION} \
            VERSION=${VERSION} \
            COMMIT=${COMMIT} \
            LEDGER_ENABLED=${LEDGER_ENABLED}
        mv ./build/${APP}${exe_file_extension} ${OUTDIR}/${APP}-${VERSION}-${os}-${arch}${exe_file_extension}
    done
    unset exe_file_extension
done

# Generate and display build report
f_generate_build_report ${OUTDIR}
cat ${OUTDIR}/build_report
```

# Makefile integration

An example of integration with the client application's `Makefile` follows:

```Makefile
VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true

build-simd-all: go.sum
	docker pull cosmossdk/rbuilder:latest
	docker rm latest-build || true
	docker run --volume=$(CURDIR):/sources:ro \
        --env TARGET_OS='darwin linux windows' \
        --env APP=simd \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build cosmossdk/rbuilder:latest
	docker cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/
```
