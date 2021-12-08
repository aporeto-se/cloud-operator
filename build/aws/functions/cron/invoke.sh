#!/bin/bash
# shellcheck disable=SC2002

funcname="cloud-operator-cron"

function main() {
  trap cleanup EXIT
  tmpfile=$(mktemp)
  which aws > /dev/null 2>&1 || { err "aws cli not found in path"; return 2; }
  aws lambda invoke --function-name $funcname "$tmpfile" --log-type Tail --query 'LogResult' --output text | base64 -d
  rc=$?
  [ $rc -eq 0 ] || return $rc
  [ -f "$tmpfile" ] || return 0
  which jq > /dev/null 2>&1 && { cat "$tmpfile" | jq; return 0; }
  cat "$tmpfile"
}

function cleanup() { [[ "$tmpfile" ]] && rm -rf "$tmpfile"; }

function err() { echo "$@" 1>&2; }

main "$@"