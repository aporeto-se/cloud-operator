#!/bin/bash -e

funcname="cloud-operator-cron"
cd "$(dirname "$0")/build"
aws lambda update-function-code --function-name "$funcname" --zip-file fileb://function.zip --no-cli-pager