# Makefile for vela-go-definitions

# Go parameters
GOCMD=go
GOMOD=$(GOCMD) mod

# Ginkgo parameters
GINKGO=$(shell which ginkgo 2>/dev/null || echo "go run github.com/onsi/ginkgo/v2/ginkgo")

# Test data path
TESTDATA_PATH ?= test/builtin-definition-example

# Generated definitions output directory
DEFINITIONS_DIR ?= vela-templates/definitions

# Timeout for E2E tests
E2E_TIMEOUT ?= 10m

# Number of parallel processes for Ginkgo (can be overridden)
PROCS ?= 10

# Baseline directory for parity tests
BASELINE_DIR ?= /tmp/cue-baseline

# k3d cluster name for local E2E testing
E2E_CLUSTER ?= e2e-test


.PHONY: tidy install-ginkgo test-unit test-e2e test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps test-e2e-parity generate-baseline e2e-setup e2e-teardown cleanup-e2e-namespaces force-cleanup-e2e-namespaces generate fmt vet lint check-diff reviewable help

## Generate CUE definitions from Go into vela-templates/definitions/
generate:
	@echo "Generating CUE definitions..."
	$(GOCMD) run ./cmd/defkit generate --output-dir $(DEFINITIONS_DIR)

## Format Go code
fmt:
	@echo "Formatting Go code..."
	$(GOCMD) fmt ./...

## Vet Go code
vet:
	@echo "Vetting Go code..."
	$(GOCMD) vet ./...

## Lint Go code (requires golangci-lint)
lint:
	@echo "Linting Go code..."
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## Check that generated files are up-to-date (no uncommitted diff after generate)
check-diff: generate
	@echo "Checking for uncommitted changes..."
	@if git diff --quiet -- $(DEFINITIONS_DIR); then \
		echo "Generated definitions are up-to-date."; \
	else \
		echo "ERROR: Generated definitions are out of date. Run 'make generate' and commit the changes."; \
		git diff --stat -- $(DEFINITIONS_DIR); \
		exit 1; \
	fi

## Run all reviewable checks: generate, format, vet, lint, check-diff
reviewable: generate fmt vet lint check-diff

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

test-e2e-components: force-cleanup-e2e-namespaces
	@echo "Running E2E tests for component definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="components" --procs=$(PROCS) ./test/e2e/...

test-e2e-traits: force-cleanup-e2e-namespaces
	@echo "Running E2E tests for trait definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="traits" --procs=$(PROCS) ./test/e2e/...

test-e2e-policies: force-cleanup-e2e-namespaces
	@echo "Running E2E tests for policy definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="policies" --procs=$(PROCS) ./test/e2e/...

test-e2e-workflowsteps: force-cleanup-e2e-namespaces
	@echo "Running E2E tests for workflowstep definitions in parallel ($(PROCS) processes)..."
	TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="workflowsteps" --procs=$(PROCS) ./test/e2e/...

## Dry-run parity tests (compares defkit output against CUE baseline)
test-e2e-parity:
	@echo "Running dry-run parity tests ($(PROCS) processes)..."
	CUE_BASELINE_DIR=$(BASELINE_DIR) TESTDATA_PATH=$(TESTDATA_PATH) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="parity" --procs=$(PROCS) ./test/e2e/...

## Generate CUE baseline for parity tests (requires built-in CUE definitions installed)
generate-baseline:
	@echo "Generating CUE dry-run baseline into $(BASELINE_DIR)..."
	@mkdir -p $(BASELINE_DIR)
	@for dir in components trait policies workflowsteps; do \
		for f in $(TESTDATA_PATH)/applications/$$dir/*.yaml; do \
			[ -f "$$f" ] || continue; \
			base=$$(basename "$$f" .yaml); \
			echo "  dry-run: $$dir/$$base"; \
			vela dry-run -f "$$f" > "$(BASELINE_DIR)/$$base.yaml" 2>/dev/null || true; \
		done; \
	done
	@echo "Baseline generated at $(BASELINE_DIR)"

## Set up a local E2E test environment (k3d cluster + KubeVela + defkit definitions)
## Prerequisites: docker, k3d, kubectl, vela CLI
e2e-setup:
	@echo "=== Setting up E2E test environment ==="
	@# Step 1: Fix Docker inotify limits (required for k3s on macOS)
	@echo "[1/6] Fixing inotify limits in Docker VM..."
	@docker run --rm --privileged alpine:latest sh -c \
		"sysctl -w fs.inotify.max_user_watches=524288 > /dev/null && sysctl -w fs.inotify.max_user_instances=512 > /dev/null" 2>/dev/null || true
	@# Step 2: Create k3d cluster
	@echo "[2/6] Creating k3d cluster '$(E2E_CLUSTER)'..."
	@k3d cluster delete $(E2E_CLUSTER) 2>/dev/null || true
	@k3d cluster create $(E2E_CLUSTER) --wait --timeout 180s
	@kubectl config use-context k3d-$(E2E_CLUSTER)
	@# Step 3: Wait for node
	@echo "[3/6] Waiting for node to be ready..."
	@for i in $$(seq 1 60); do \
		nodes=$$(kubectl get nodes --no-headers 2>/dev/null | wc -l | tr -d ' '); \
		if [ "$$nodes" -gt 0 ]; then \
			kubectl get nodes; \
			break; \
		fi; \
		if [ "$$i" -eq 60 ]; then echo "ERROR: Node did not become ready"; exit 1; fi; \
		sleep 5; \
	done
	@# Step 4: Install KubeVela
	@echo "[4/6] Installing KubeVela..."
	@vela install
	@kubectl wait --for=condition=available --timeout=300s deployment/kubevela-vela-core -n vela-system
	@# Step 5: Uninstall built-in definitions and install defkit definitions
	@echo "[5/6] Replacing built-in definitions with defkit definitions..."
	@kubectl delete componentdefinitions --all -n vela-system 2>/dev/null || true
	@kubectl delete traitdefinitions --all -n vela-system 2>/dev/null || true
	@kubectl delete workflowstepdefinitions --all -n vela-system 2>/dev/null || true
	@kubectl delete policydefinitions --all -n vela-system 2>/dev/null || true
	@$(MAKE) generate
	@for dir in $(DEFINITIONS_DIR)/*/; do \
		for f in $$dir*.cue; do \
			[ -f "$$f" ] && vela def apply "$$f" 2>&1; \
		done; \
	done
	@# Step 6: Install ginkgo
	@echo "[6/6] Installing Ginkgo..."
	@$(MAKE) install-ginkgo
	@echo ""
	@echo "=== E2E environment ready ==="
	@echo "Cluster:      k3d-$(E2E_CLUSTER)"
	@echo "Definitions:  $$(kubectl get componentdefinitions,traitdefinitions,policydefinitions,workflowstepdefinitions -n vela-system --no-headers 2>/dev/null | wc -l | tr -d ' ') installed"
	@echo ""
	@echo "Run tests with:"
	@echo "  make test-e2e-components"
	@echo "  make test-e2e-traits"
	@echo "  make test-e2e-policies"
	@echo "  make test-e2e-workflowsteps"
	@echo "  make test-e2e                  # all of the above"

## Tear down the local E2E test environment
e2e-teardown:
	@echo "Tearing down E2E test environment..."
	@k3d cluster delete $(E2E_CLUSTER) 2>/dev/null || true
	@echo "Cluster '$(E2E_CLUSTER)' deleted."

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
	@echo ""
	@echo "  Reviewable:"
	@echo "  reviewable             - Run all checks: generate, fmt, vet, lint, check-diff"
	@echo "  generate               - Generate CUE definitions from Go into vela-templates/definitions/"
	@echo "  fmt                    - Format Go code"
	@echo "  vet                    - Vet Go code"
	@echo "  lint                   - Lint Go code (installs golangci-lint if missing)"
	@echo "  check-diff             - Verify generated definitions are up-to-date"
	@echo ""
	@echo "  Dependencies:"
	@echo "  tidy                   - Tidy go.mod dependencies"
	@echo "  install-ginkgo         - Install Ginkgo CLI for running E2E tests"
	@echo ""
	@echo "  Tests:"
	@echo "  test-unit              - Run unit tests (no cluster required)"
	@echo "  test-e2e               - Run all E2E tests"
	@echo "  test-e2e-components    - Run E2E tests for component definitions (parallel)"
	@echo "  test-e2e-traits        - Run E2E tests for trait definitions (parallel)"
	@echo "  test-e2e-policies      - Run E2E tests for policy definitions (parallel)"
	@echo "  test-e2e-workflowsteps - Run E2E tests for workflowstep definitions (parallel)"
	@echo "  test-e2e-parity        - Run dry-run parity tests (defkit vs CUE baseline)"
	@echo "  generate-baseline      - Generate CUE dry-run baseline for parity tests"
	@echo ""
	@echo "  Environment:"
	@echo "  e2e-setup                    - Set up local E2E environment (k3d + KubeVela + definitions)"
	@echo "  e2e-teardown                 - Tear down local E2E environment"
	@echo ""
	@echo "  Cleanup:"
	@echo "  cleanup-e2e-namespaces       - Delete all namespaces starting with 'e2e'"
	@echo "  force-cleanup-e2e-namespaces - Force delete stuck terminating namespaces starting with 'e2e'"
	@echo ""
	@echo "Environment variables:"
	@echo "  DEFINITIONS_DIR - Output directory for generated CUE (default: vela-templates/definitions)"
	@echo "  TESTDATA_PATH   - Path to test data (default: test/builtin-definition-example)"
	@echo "  E2E_TIMEOUT     - Timeout for E2E tests (default: 10m)"
	@echo "  PROCS           - Number of parallel processes for Ginkgo (default: 10)"
	@echo "  BASELINE_DIR    - Directory for CUE parity baselines (default: /tmp/cue-baseline)"
	@echo ""
	@echo "Examples:"
	@echo "  make e2e-setup                              # Set up local test cluster"
	@echo "  make test-e2e                               # Run all E2E tests"
	@echo "  make test-e2e-components PROCS=4            # Run component tests with 4 processes"
	@echo "  make e2e-teardown                           # Tear down test cluster"
	@echo "  make reviewable                             # Run all pre-submit checks"

