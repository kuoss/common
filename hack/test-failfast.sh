#!/bin/bash
set -euo pipefail
cd $(dirname $0)/../
for s in $(go list ./...); do
    echo =============== $s ===============
    go test -race -failfast -v $s
done