COVERAGE_THRESHOLD = 90

.PHONY: test
test:
	go test -race --failfast ./...

#### checks
.PHONY: checks
checks: cover lint licenses vulncheck

.PHONY: cover
cover: test
	@echo "Running tests and checking coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' | \
	awk -v threshold=$(COVERAGE_THRESHOLD) '{ if ($$1 < threshold) { print "Coverage is below threshold: " $$1 "% (<" threshold "%)"; exit 1 } else { print "Coverage is sufficient: " $$1 "% (>=" threshold "%)" } }'

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run

.PHONY: licenses
licenses: go-licenses
	$(GO_LICENSES) check ./...

.PHONY: vulncheck
vulncheck: govulncheck
	$(GOVULNCHECK) ./...

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOVULNCHECK ?= $(LOCALBIN)/govulncheck
GO_LICENSES ?= $(LOCALBIN)/go-licenses

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.60.2
GOVULNCHECK_VERSION ?= latest
GO_LICENSES_VERSION ?= v1.6.0

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: govulncheck
govulncheck: $(GOVULNCHECK)
$(GOVULNCHECK): $(LOCALBIN)
	$(call go-install-tool,$(GOVULNCHECK),golang.org/x/vuln/cmd/govulncheck,$(GOVULNCHECK_VERSION))

.PHONY: go-licenses
go-licenses: $(GO_LICENSES)
$(GO_LICENSES): $(LOCALBIN)
	$(call go-install-tool,$(GO_LICENSES),github.com/google/go-licenses,$(GO_LICENSES_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
