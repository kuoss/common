#!/bin/bash
cd $(dirname $0)/../

which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --timeout 5m
