#!/bin/bash

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -o pipefail

if [ -z "${OS:-$(go env GOOS)}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-$(go env GOARCH)}" ]; then
    echo "ARCH must be set"
    exit 1
fi

# Stamp the current date as the version.
VERSION="v$(date +%F)"
sed -i \
  "s/\(const Version =\).*$/\1 \"${VERSION}\"/" \
  main.go

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

# Create a file called `s3cli` in the current build directory.

go mod tidy
go mod vendor

gofmt -w .

GOBIN="$(pwd)" go install   \
    -installsuffix "static" \
    -ldflags "-s -w"        \
    ./...
