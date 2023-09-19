#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go mod tidy
go fmt ./...
go vet ./...

./hack/misspell.sh
./hack/gocyclo.sh

./hack/goimports.sh
./hack/staticcheck.sh
./hack/golangci-lint.sh

./hack/test-cover.sh
