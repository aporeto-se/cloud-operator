#!/bin/bash -e

function main() {
  cd "$(dirname "$0")"

  out=$(./linux-amd64) || {
    err "Failed"
    return 3
  }
  which jq > /dev/null 2>&1 && {
    err "Sending output to jq"
    echo "$out" | jq
    return 0
  }
  err "Sending output to stdout"
  echo "$out"
  return 0
}

function err() { echo "$@" 1>&2; }

main "$@"