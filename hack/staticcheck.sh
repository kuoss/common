#!/bin/bash
cd $(dirname $0)/../

which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
