#!/bin/bash

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

_gofmt() {
  test -z "$(gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") | tee /dev/stderr)"
}

_go_test() {
  for dir in $(_test_dirs);
  do
    go test -v "${dir}"
  done
}

_test_dirs() {
  local test_dirs=$(find . -type f -name '*.go' -not -path "./.*" -not -path "*/_*" -not -path "./Godeps/*" -not -path "./vendor/*" -print0 | xargs -0 -n1 dirname| sort -u)
  echo "${test_dirs}"
}


test_unit() {
  local go_tests
  go_tests=(gofmt go_test)

  echo "available unit tests: ${go_tests[*]}"

  ((n_elements=${#go_tests[@]}, max=n_elements - 1))

  for ((i = 0; i <= max; i++)); do
     echo "running ${go_tests[i]}"
      _"${go_tests[i]}"
  done
}

test_unit