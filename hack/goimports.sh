#!/bin/bash
cd $(dirname $0)/../

which goimports || go install golang.org/x/tools/cmd/goimports@latest
goimports -local -v -w .
