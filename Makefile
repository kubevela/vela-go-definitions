# Makefile for vela-go-definitions

# Go parameters
GOCMD=go
GOMOD=$(GOCMD) mod

# Ginkgo parameters
GINKGO=$(shell which ginkgo 2>/dev/null || echo "go run github.com/onsi/ginkgo/v2/ginkgo")

# Test data path
TESTDATA_PATH ?= test/builtin-definition-example


# Timeout for E2E tests
E2E_TIMEOUT ?= 10m

# Number of parallel processes for Ginkgo (can be overridden)
PROCS ?= 10

.PHONY: tidy install-ginkgo test-unit test-e2e test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps cleanup-e2e-namespaces force-cleanup-e2e-namespaces help

## Dependency management
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## Install Ginkgo CLI
install-ginkgo:
	@echo "Installing Ginkgo CLI..."
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

## Unit tests
test-unit:
	@echo "Running unit tests..."
	$(GOCMD) test -v -race -count=1 ./components/... ./traits/... ./policies/... ./workflowsteps/...

## E2E Test targets
test-e2e: test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps
	@echo "All E2E tests completed!"

test-e2e-components:
	@echo "Running E2E tests for component definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="components" --procs=$(PROCS) ./test/e2e/...

test-e2e-traits:
	@echo "Running E2E tests for trait definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="traits" --procs=$(PROCS) ./test/e2e/...

test-e2e-policies:
	@echo "Running E2E tests for policy definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="policies" --procs=$(PROCS) ./test/e2e/...

test-e2e-workflowsteps:
	@echo "Running E2E tests for workflowstep definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="workflowsteps" --procs=$(PROCS) ./test/e2e/...

## Cleanup E2E test namespaces
cleanup-e2e-namespaces:
	@echo "Deleting all namespaces starting with 'e2e'..."
	@kubectl get namespaces --no-headers -o custom-columns=":metadata.name" | grep "^e2e" | xargs -r kubectl delete namespace --wait=false || true
	@echo "Cleanup complete!"

## Force cleanup E2E test namespaces (removes finalizers for stuck namespaces)
force-cleanup-e2e-namespaces:
	@echo "Force deleting all namespaces starting with 'e2e'..."
	@for ns in $$(kubectl get namespaces --no-headers -o custom-columns=":metadata.name" | grep "^e2e"); do \
		echo "Force deleting namespace: $$ns"; \
		kubectl get namespace $$ns -o json | jq '.spec.finalizers = []' | kubectl replace --raw "/api/v1/namespaces/$$ns/finalize" -f - || true; \
	done
	@echo "Force cleanup complete!"

## Help
help:
	@echo "Available targets:"
	@echo "  tidy                   - Tidy go.mod dependencies"
	@echo "  install-ginkgo         - Install Ginkgo CLI for running E2E tests"
	@echo "  test-unit              - Run unit tests (no cluster required)"
	@echo "  test-e2e               - Run all E2E tests"
	@echo "  test-e2e-components    - Run E2E tests for component definitions (parallel)"
	@echo "  test-e2e-traits        - Run E2E tests for trait definitions (parallel)"
	@echo "  test-e2e-policies      - Run E2E tests for policy definitions (parallel)"
	@echo "  test-e2e-workflowsteps - Run E2E tests for workflowstep definitions (parallel)"
	@echo "  cleanup-e2e-namespaces       - Delete all namespaces starting with 'e2e'"
	@echo "  force-cleanup-e2e-namespaces - Force delete stuck terminating namespaces starting with 'e2e'"
	@echo ""
	@echo "Environment variables:"
	@echo "  TESTDATA_PATH - Path to test data (default: test/builtin-definition-example)"
	@echo "  E2E_TIMEOUT   - Timeout for E2E tests (default: 30m)"
	@echo "  PROCS         - Number of parallel processes for Ginkgo (default: 4)"
	@echo ""
	@echo "Examples:"
	@echo "  make test-e2e-components                    # Run with 4 processes (default)"
	@echo "  make test-e2e-components PROCS=8            # Run with 8 processes"
	@echo "  make test-e2e-components PROCS=2            # Run with 2 processes"

