# Makefile for vela-go-definitions

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Ginkgo parameters
GINKGO=$(shell which ginkgo 2>/dev/null || echo "go run github.com/onsi/ginkgo/v2/ginkgo")

# Test data path
TESTDATA_PATH ?= test/builtin-definition-example

# Vela CLI path (can be overridden)
VELA_CLI ?= vela

# Timeout for E2E tests
E2E_TIMEOUT ?= 30m

.PHONY: all build test test-unit test-e2e test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps clean help

all: build test

## Build targets
build:
	@echo "Building..."
	$(GOBUILD) ./...

## Test targets
test: test-unit
	@echo "All unit tests passed!"

test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -v ./components/... ./traits/... ./policies/... ./workflowsteps/...

## E2E Test targets
test-e2e: test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps
	@echo "All E2E tests completed!"

test-e2e-components:
	@echo "Running E2E tests for component definitions..."
	@echo "Test file: test/e2e/component_e2e_test.go"
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="components" ./test/e2e/...

test-e2e-traits:
	@echo "Running E2E tests for trait definitions..."
	@echo "Test file: test/e2e/trait_e2e_test.go"
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="traits" ./test/e2e/...

test-e2e-policies:
	@echo "Running E2E tests for policy definitions..."
	@echo "Test file: test/e2e/policy_e2e_test.go"
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="policies" ./test/e2e/...

test-e2e-workflowsteps:
	@echo "Running E2E tests for workflowstep definitions..."
	@echo "Test file: test/e2e/workflowstep_e2e_test.go"
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="workflowsteps" ./test/e2e/...

## Install Ginkgo CLI
install-ginkgo:
	@echo "Installing Ginkgo CLI..."
	$(GOGET) github.com/onsi/ginkgo/v2/ginkgo@latest

## Dependency management
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## Clean targets
clean:
	@echo "Cleaning..."
	rm -rf bin/
	$(GOCMD) clean -testcache

## Help
help:
	@echo "Available targets:"
	@echo "  all                    - Build and run unit tests"
	@echo "  build                  - Build all packages"
	@echo "  test                   - Run unit tests"
	@echo "  test-unit              - Run unit tests"
	@echo "  test-e2e               - Run all E2E tests"
	@echo "  test-e2e-components    - Run E2E tests for component definitions"
	@echo "                           (test/e2e/component_e2e_test.go)"
	@echo "  test-e2e-traits        - Run E2E tests for trait definitions"
	@echo "                           (test/e2e/trait_e2e_test.go)"
	@echo "  test-e2e-policies      - Run E2E tests for policy definitions"
	@echo "                           (test/e2e/policy_e2e_test.go)"
	@echo "  test-e2e-workflowsteps - Run E2E tests for workflowstep definitions"
	@echo "                           (test/e2e/workflowstep_e2e_test.go)"
	@echo "  install-ginkgo         - Install Ginkgo CLI"
	@echo "  deps                   - Download dependencies"
	@echo "  tidy                   - Tidy go.mod"
	@echo "  clean                  - Clean build artifacts"
	@echo ""
	@echo "Test files:"
	@echo "  test/e2e/helpers_test.go          - Common helper functions"
	@echo "  test/e2e/component_e2e_test.go    - Component definition tests"
	@echo "  test/e2e/trait_e2e_test.go        - Trait definition tests"
	@echo "  test/e2e/policy_e2e_test.go       - Policy definition tests"
	@echo "  test/e2e/workflowstep_e2e_test.go - WorkflowStep definition tests"
	@echo ""
	@echo "Environment variables:"
	@echo "  TESTDATA_PATH - Path to test data (default: test/builtin-definition-example)"
	@echo "  VELA_CLI      - Path to vela CLI (default: vela)"
	@echo "  E2E_TIMEOUT   - Timeout for E2E tests (default: 30m)"

