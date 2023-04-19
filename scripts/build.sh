#!/usr/bin/env bash

set -ex

GO=${GO:-$(which go)}
GOOS=${GOOS:-$(${GO} env GOOS)}
GOARCH=${GOARCH:amd64}

#LDFLAGS="-extldflags -static"
LDFLAGS=""

BUILD_ENTRY=${BUILD_ENTRY:-$1}
OUT=${OUT:-$2}

CGO_ENABLED=0 GOARCH="${GOARCH}" GOOS=${GOOS} \
  "${GO}" build -o "${OUT}" -ldflags "${LDFLAGS}" "${BUILD_ENTRY}"

