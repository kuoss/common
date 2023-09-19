
test:
	hack/test-failfast.sh

cover:
	hack/test-cover.sh

checks:
	hack/checks.sh


go-licenses:
	hack/go-licenses.sh
gocyclo:
	hack/gocyclo.sh
goimports:
	hack/goimports.sh
golangci-lint:
	hack/golangci-lint.sh
misspell:
	hack/misspell.sh
staticcheck:
	hack/staticcheck.sh
