#!/bin/bash

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"


name=${__proj_dir##*/}
build_dir="${__proj_dir}/build"
go_build=(go build -ldflags "-w")

export CGO_ENABLED=0

# rebuild binaries:
echo "removing: ${build_dir:?}/*"
rm -rf "${build_dir:?}/"*

echo "building numbers"
export GOOS=linux
export GOARCH=amd64
mkdir -p "${build_dir}/${GOOS}/x86_64"
"${go_build[@]}" -o "${build_dir}/${GOOS}/x86_64/${name}" . || exit 1