
test:
	hack/test-failfast.sh

cover:
	hack/test-cover.sh

checks:
	hack/checks.sh

misspell:
	hack/misspell.sh

gocyclo:
	hack/gocyclo.sh