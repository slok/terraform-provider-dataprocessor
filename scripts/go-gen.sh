#!/bin/bash
# vim: ai:ts=8:sw=8:noet
set -efCo pipefail
export SHELLOPTS
IFS=$'\t\n'

command -v "go" >/dev/null 2>&1 || {
    echo "please install go"
    exit 1
}

go generate ./...
