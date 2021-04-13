VENDORDIR = vendor

GO ?= go
GOLANGCI_LINT ?= golangci-lint

.PHONY: $(VENDORDIR) lint test test-unit

$(VENDORDIR):
	@mkdir -p $(VENDORDIR)
	@$(GO) mod vendor

lint:
	@$(GOLANGCI_LINT) run

test: test-unit test-integration

## Run unit tests
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

test-integration:
	@echo ">> integration test"
	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./...  -godog -race
