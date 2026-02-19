# E2E Test Automation Framework

This document provides comprehensive documentation of the E2E test automation framework implemented in the `feat/defkit-test-automation` branch for the `vela-go-definitions` module.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Directory Structure](#directory-structure)
- [Test Framework](#test-framework)
  - [Ginkgo Test Suite](#ginkgo-test-suite)
  - [Test Helpers](#test-helpers)
  - [Test Data](#test-data)
- [CI/CD Integration](#cicd-integration)
  - [GitHub Actions Workflows](#github-actions-workflows)
  - [Reusable Setup Action](#reusable-setup-action)
- [Running Tests Locally](#running-tests-locally)
- [How It Works](#how-it-works)
- [Adding New Tests](#adding-new-tests)
- [Troubleshooting](#troubleshooting)

---

## Overview

The E2E test automation framework validates that KubeVela X-Definitions (components, traits, policies, and workflow steps) defined in this module work correctly when deployed to a live Kubernetes cluster with KubeVela installed.

### Key Features

- **Automated E2E Testing**: Validates definitions against a real KubeVela cluster
- **Parallel Execution**: Tests run in parallel using Ginkgo's multi-process support
- **Isolated Test Namespaces**: Each test creates a unique namespace to avoid conflicts
- **CI/CD Integration**: GitHub Actions workflows automatically test on PR/push
- **Dynamic KubeVela Version**: Builds KubeVela CLI from the commit specified in `go.mod`

---

## Architecture

### High-Level Overview

The diagram below shows the complete CI/CD pipeline flow:

1. **Trigger Events** - Workflows are triggered by git push, pull requests (when `go.mod`/`go.sum` change), or manual dispatch
2. **GitHub Actions Workflows** - Four separate workflows run tests for components, traits, policies, and workflow steps
3. **Setup Action** - A reusable composite action sets up the entire test environment (Kind cluster, KubeVela, definitions)
4. **Test Execution** - Ginkgo runs E2E tests in parallel, each test in its own isolated namespace
5. **Results** - Test results are reported back to GitHub with detailed summaries

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                              TRIGGER EVENTS                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌──────────────────┐                      │
│  │  git push   │    │ Pull Request│    │ workflow_dispatch│                      │
│  │ (go.mod/    │    │ (go.mod/    │    │ (manual trigger) │                      │
│  │  go.sum)    │    │  go.sum)    │    │                  │                      │
│  └──────┬──────┘    └──────┬──────┘    └────────┬─────────┘                      │
│         └──────────────────┼───────────────────┘                                 │
│                            ▼                                                     │
│  ┌────────────────────────────────────────────────────────────────────────────┐  │
│  │                        GitHub Actions Workflows                            │  │
│  │  ┌────────────────┐ ┌────────────────┐ ┌────────────────┐ ┌──────────────┐ │  │
│  │  │test-component- │ │test-trait-     │ │test-policy-    │ │test-workflow-│ │  │
│  │  │definitions.yaml│ │definitions.yaml│ │definitions.yaml│ │step-def.yaml │ │  │
│  │  └───────┬────────┘ └───────┬────────┘ └───────┬────────┘ └──────┬───────┘ │  │
│  │          └──────────────────┼──────────────────┼─────────────────┘         │  │
│  │                             ▼                  ▼                           │  │
│  │              ┌──────────────────────────────────────┐                      │  │
│  │              │  setup-vela-environment (Composite)  │                      │  │
│  │              │  Sets up: Go, Kind, KubeVela, Defs   │                      │  │
│  │              └──────────────┬───────────────────────┘                      │  │
│  │                             │                                              │  │
│  │                             ▼                                              │  │
│  │              ┌──────────────────────────────────────┐                      │  │
│  │              │  make test-e2e-{components|traits|   │                      │  │
│  │              │       policies|workflowsteps}        │                      │  │
│  │              │  Runs Ginkgo E2E tests in parallel   │                      │  │
│  │              └──────────────┬───────────────────────┘                      │  │
│  │                             │                                              │  │
│  │                             ▼                                              │  │
│  │              ┌──────────────────────────────────────┐                      │  │
│  │              │         Test Results Summary         │                      │  │
│  │              │  ✅ Passed | ❌ Failed | ⏭ Skipped  │                      │  │
│  │              └──────────────────────────────────────┘                      │  │
│  └────────────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

### Setup Action Detail Flow

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                    setup-vela-environment (Composite Action)                     │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  PHASE 1: Environment Setup                                                      │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                             │ │
│  │   ┌─────────────┐      ┌──────────────────┐      ┌───────────────────┐      │ │
│  │   │  Setup Go   │ ───► │ Create Kind      │ ───► │ Verify Cluster    │      │ │
│  │   │  (v1.23)    │      │ Cluster          │      │ kubectl get nodes │      │ │
│  │   └─────────────┘      │ (vela-test)      │      └───────────────────┘      │ │
│  │                        └──────────────────┘                                 │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  PHASE 2: KubeVela Installation                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                             │ │
│  │   ┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐    │ │
│  │   │ Download Vela   │ ───► │ vela install     │ ───► │ Wait for        │    │ │
│  │   │ CLI (latest)    │      │ (deploy to       │      │ controller ready│    │ │
│  │   │                 │      │  cluster)        │      │ (300s timeout)  │    │ │
│  │   └─────────────────┘      └──────────────────┘      └─────────────────┘    │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  PHASE 3: Build CLI from go.mod Reference                                        │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                             │ │
│  │   ┌─────────────────────────────────────────────────────────────────────┐   │ │
│  │   │                     Parse go.mod replace directive                  │   │ │
│  │   │  ┌─────────────────────────────────────────────────────────────┐    │   │ │
│  │   │  │ replace github.com/oam-dev/kubevela =>                      │    │   │ │
│  │   │  │         github.com/kubevela/kubevela v0.0.0-20250115-abc123 │    │   │ │
│  │   │  └─────────────────────────────────────────────────────────────┘    │   │ │
│  │   │                              │                                      │   │ │
│  │   │              ┌───────────────┴───────────────┐                      │   │ │
│  │   │              ▼                               ▼                      │   │ │
│  │   │   ┌──────────────────┐            ┌──────────────────┐              │   │ │
│  │   │   │ Extract repo:    │            │ Extract commit:  │              │   │ │
│  │   │   │ kubevela/kubevela│            │ abc123           │              │   │ │
│  │   │   └──────────────────┘            └──────────────────┘              │   │ │
│  │   └─────────────────────────────────────────────────────────────────────┘   │ │
│  │                              │                                              │ │
│  │                              ▼                                              │ │
│  │   ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐   │ │
│  │   │ git clone         │───►│ git checkout      │───►│ make vela-cli     │   │ │
│  │   │ kubevela/kubevela │    │ abc123            │    │ (build from src)  │   │ │
│  │   └───────────────────┘    └───────────────────┘    └───────────────────┘   │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  PHASE 4: Install Definitions from This Module                                   │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                             │ │
│  │   ┌───────────────────────┐    ┌───────────────────────────────────────┐    │ │
│  │   │ Uninstall built-in    │───►│ vela def apply-module .               │    │ │
│  │   │ definitions           │    │ --conflict=overwrite                  │    │ │
│  │   │ (componentdefs,       │    │                                       │    │ │
│  │   │  traitdefs, etc.)     │    │ Installs definitions from:            │    │ │
│  │   └───────────────────────┘    │  - components/*.go                    │    │ │
│  │                                │  - traits/*.go                        │    │ │
│  │                                │  - policies/*.go                      │    │ │
│  │                                │  - workflowsteps/*.go                 │    │ │
│  │                                └───────────────────────────────────────┘    │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  PHASE 5: Prepare Test Environment                                               │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │   ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐   │ │
│  │   │ make tidy         │───►│ Verify installed  │───►│ Install Ginkgo    │   │ │
│  │   │ (go mod tidy)     │    │ definitions       │    │ CLI               │   │ │
│  │   └───────────────────┘    └───────────────────┘    └───────────────────┘   │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  PHASE 6: Run E2E Tests                                                          │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                                                                             │ │
│  │   ┌─────────────────────────────────────────────────────────────────────┐   │ │
│  │   │  make test-e2e-{components|traits|policies|workflowsteps}           │   │ │
│  │   │                                                                     │   │ │
│  │   │  Executes:                                                          │   │ │
│  │   │  ginkgo -v --timeout=30m --label-filter="<type>" --procs=4          │   │ │
│  │   │          ./test/e2e/...                                             │   │ │
│  │   └─────────────────────────────────────────────────────────────────────┘   │ │
│  │                              │                                              │ │
│  │                              ▼                                              │ │
│  │   ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐   │ │
│  │   │ For each YAML in  │───►│ Create isolated   │───►│ Apply Application │   │ │
│  │   │ test data folder  │    │ namespace (e2e-*) │    │ & wait for status │   │ │
│  │   └───────────────────┘    └───────────────────┘    └───────────────────┘   │ │
│  │                              │                                              │ │
│  │                              ▼                                              │ │
│  │   ┌─────────────────────────────────────────────────────────────────────┐   │ │
│  │   │                      Test Results                                   │   │ │
│  │   │  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐              │   │ │
│  │   │  │ ✅ PASSED   │    │ ❌ FAILED   │    │ ⏭ SKIPPED  │              │   │ │
│  │   │  │ App running │    │ Timeout/Err │    │ Filtered    │              │   │ │
│  │   │  └─────────────┘    └─────────────┘    └─────────────┘              │   │ │
│  │   └─────────────────────────────────────────────────────────────────────┘   │ │
│  │                                                                             │ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────────┘
```


### Single Test Execution Flow

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                        Single Test Execution Flow                                │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│   Input: webservice.yaml                                                         │
│   ┌─────────────────────────────────────────────────────────────────┐            │
│   │ apiVersion: core.oam.dev/v1beta1                                │            │
│   │ kind: Application                                               │            │
│   │ metadata:                                                       │            │
│   │   name: website                                                 │            │
│   │ spec:                                                           │            │
│   │   components:                                                   │            │
│   │     - name: frontend                                            │            │
│   │       type: webservice                                          │            │
│   │       properties:                                               │            │
│   │         image: oamdev/testapp:v1                                │            │
│   └─────────────────────────────────────────────────────────────────┘            │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 1: Read Application                                                  │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ app, err := readAppFromFile("webservice.yaml")                      │   │   │
│  │ │ // Parses YAML, handles multi-document files                        │   │   │
│  │ │ // Returns: Application{Name: "website", ...}                       │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 2: Generate Unique Namespace                                         │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ uniqueNs := uniqueNamespaceForApp(app.Name)                         │   │   │
│  │ │ // "website" → "e2e-website"                                        │   │   │
│  │ │ // Sanitizes: lowercase, replace . and _, max 30 chars              │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 3: Create Namespace                                                  │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ ensureNamespace(ctx, "e2e-website")                                 │   │   │
│  │ │ // kubectl create namespace e2e-website                             │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 4: Cleanup Existing Application (if any)                             │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ cleanupExistingApplication(ctx, app)                                │   │   │
│  │ │ // Deletes app if exists, waits for complete deletion (30s timeout) │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 5: Apply Prerequisites (if any)                                      │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ applyPrerequisitesIfAny(ctx, filePath, uniqueNs)                    │   │   │
│  │ │ // For ref-objects.yaml: creates Deployment/Service first           │   │   │
│  │ │ // Waits 2s for resources to settle                                 │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 6: Apply Application                                                 │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ k8sClient.Create(ctx, app)                                          │   │   │
│  │ │ // Creates Application CR in Kubernetes                             │   │   │
│  │ │ // KubeVela controller starts reconciliation                        │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 7: Wait for Running Status                                           │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ waitForApplicationRunning(ctx, "website", "e2e-website")            │   │   │
│  │ │                                                                     │   │   │
│  │ │    ┌─────────────────────────────────────────────────────────────┐  │   │   │
│  │ │    │  Poll Loop (every 5s, timeout 5m):                          │  │   │   │
│  │ │    │  ┌─────────────────────────────────────────────────────┐    │  │   │   │
│  │ │    │  │ status := getApplicationStatus(ctx, app, ns)        │    │  │   │   │
│  │ │    │  │ if status contains "running" → SUCCESS ✅           │    │  │   │   │
│  │ │    │  │ if status contains "failed"  → FAIL ❌              │    │  │   │   │
│  │ │    │  │ else → continue polling                             │    │  │   │   │
│  │ │    │  └─────────────────────────────────────────────────────┘    │  │   │   │
│  │ │    └─────────────────────────────────────────────────────────────┘  │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│  ┌───────────────────────────────────────────────────────────────────────────┐   │
│  │ STEP 8: Cleanup                                                           │   │
│  │ ┌─────────────────────────────────────────────────────────────────────┐   │   │
│  │ │ deleteNamespace(ctx, "e2e-website")                                 │   │   │
│  │ │ // Deletes entire namespace, cascading to all resources             │   │   │
│  │ └─────────────────────────────────────────────────────────────────────┘   │   │
│  └───────────────────────────────────────────────────────────────────────────┘   │
│                                      │                                           │
│                                      ▼                                           │
│                          ┌─────────────────────┐                                 │
│                          │   TEST COMPLETE ✅  │                                 │
│                          └─────────────────────┘                                 │
└──────────────────────────────────────────────────────────────────────────────────┘
```

### Parallel Execution Model

```
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         Parallel Execution (PROCS=4)                             │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│   Ginkgo distributes tests across 4 parallel processes:                          │
│                                                                                  │
│   ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│   │   Process 1     │  │   Process 2     │  │   Process 3     │  │  Process 4  │ │
│   ├─────────────────┤  ├─────────────────┤  ├─────────────────┤  ├─────────────┤ │
│   │ webservice.yaml │  │ worker.yaml     │  │ task.yaml       │  │ daemon.yaml │ │
│   │ Namespace:      │  │ Namespace:      │  │ Namespace:      │  │ Namespace:  │ │
│   │ e2e-website     │  │ e2e-vela-app    │  │ e2e-task-app    │  │ e2e-daemon  │ │
│   │                 │  │                 │  │                 │  │             │ │
│   │ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────┐ │ │
│   │ │ Deployment  │ │  │ │ Deployment  │ │  │ │ Job         │ │  │ │DaemonSet│ │ │
│   │ │ frontend    │ │  │ │ backend     │ │  │ │ task-runner │ │  │ │ daemon  │ │ │
│   │ └─────────────┘ │  │ └─────────────┘ │  │ └─────────────┘ │  │ └─────────┘ │ │
│   └────────┬────────┘  └────────┬────────┘  └────────┬────────┘  └──────┬──────┘ │
│            │                    │                    │                  │        │
│            ▼                    ▼                    ▼                  ▼        │
│   ┌─────────────────────────────────────────────────────────────────────────────┐│
│   │                        Kubernetes Cluster (Kind)                            ││
│   │  ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐   │││
│   │  │ ns:e2e-website│ │ns:e2e-vela-app│ │ns:e2e-task-app│ │ ns:e2e-daemon │   │││
│   │  │   (isolated)  │ │   (isolated)  │ │   (isolated)  │ │   (isolated)  │   │││
│   │  └───────────────┘ └───────────────┘ └───────────────┘ └───────────────┘   │││
│   └─────────────────────────────────────────────────────────────────────────────┘│
│                                                                                  │
│   Benefits:                                                                      │
│   ✓ Tests don't interfere with each other (isolated namespaces)                  │
│   ✓ 4x faster than serial execution                                              │
│   ✓ Each process has independent K8s client                                      │
│   ✓ Failures in one process don't affect others                                  │
│                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
vela-go-definitions/
├── .github/
│   ├── actions/
│   │   └── setup-vela-environment/
│   │       └── action.yaml              # Reusable composite action
│   └── workflows/
│       ├── test-component-definitions.yaml
│       ├── test-trait-definitions.yaml
│       ├── test-policy-definitions.yaml
│       └── test-workflowstep-definitions.yaml
├── test/
│   ├── e2e/
│   │   ├── e2e_suite_test.go           # Ginkgo suite setup
│   │   ├── helpers_test.go             # Shared test utilities
│   │   ├── component_e2e_test.go       # Component definition tests
│   │   ├── trait_e2e_test.go           # Trait definition tests
│   │   ├── policy_e2e_test.go          # Policy definition tests
│   │   └── workflowstep_e2e_test.go    # Workflow step tests
│   └── builtin-definition-example/
│       ├── components/                  # Component test applications
│       │   ├── webservice.yaml
│       │   ├── worker.yaml
│       │   ├── task.yaml
│       │   ├── cron-task.yaml
│       │   ├── daemon.yaml
│       │   ├── statefulset.yaml
│       │   ├── k8s-objects.yaml
│       │   └── ref-objects.yaml
│       ├── trait/                       # Trait test applications
│       │   ├── sidecar.yaml
│       │   ├── scaler.yaml
│       │   ├── gateway.yaml
│       │   ├── hpa.yaml
│       │   └── ... (28 total)
│       ├── policies/                    # Policy test applications
│       │   ├── topology.yaml
│       │   ├── override.yaml
│       │   ├── garbage-collect.yaml
│       │   └── ... (9 total)
│       └── workflowsteps/               # Workflow step test applications
│           ├── apply-component.yaml
│           ├── suspend.yaml
│           ├── deploy.yaml
│           └── ... (34 total)
└── Makefile                             # Make targets for running tests
```

---

## Test Framework

### Ginkgo Test Suite

The E2E tests use [Ginkgo v2](https://onsi.github.io/ginkgo/) with [Gomega](https://onsi.github.io/gomega/) matchers.

#### Suite Setup (`e2e_suite_test.go`)

```go
func TestE2E(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "E2E Definition Test Suite")
}

var _ = BeforeSuite(func() {
    err := initK8sClient()
    Expect(err).NotTo(HaveOccurred(), "Failed to initialize K8s client")
})
```

#### Test Structure (`component_e2e_test.go`)

Each test file follows this pattern:

```go
var _ = Describe("Component Definition E2E Tests", Label("components"), func() {
    ctx := context.Background()

    Context("when testing component definitions", func() {
        testDataPath := filepath.Join(getTestDataPath(), "components")
        componentFiles := mustListYAMLFiles(testDataPath)

        When("applying component applications", func() {
            for _, file := range componentFiles {
                file := file  // Capture for closure

                It(fmt.Sprintf("should run %s", filepath.Base(file)), func() {
                    runComponentTest(ctx, file)
                })
            }
        })
    })
})
```

#### Test Labels

Tests are organized by labels for selective execution:

| Label | Description |
|-------|-------------|
| `components` | Component definition tests |
| `traits` | Trait definition tests |
| `policies` | Policy definition tests |
| `workflowsteps` | Workflow step definition tests |

### Test Helpers

The `helpers_test.go` file provides shared utilities:

#### Key Functions

| Function | Description |
|----------|-------------|
| `initK8sClient()` | Initializes the Kubernetes controller-runtime client |
| `readAppFromFile(filename)` | Reads an Application from a YAML file (supports multi-doc) |
| `sanitizeForNamespace(name)` | Creates DNS-1123 compliant namespace names |
| `applyApplication(ctx, app)` | Creates or updates a KubeVela Application |
| `deleteApplication(ctx, app)` | Deletes an Application with foreground deletion |
| `waitForApplicationRunning(ctx, appName, namespace)` | Polls until app reaches running state |
| `hasPrerequisiteResources(filePath)` | Checks if YAML contains non-Application resources |
| `applyPrerequisiteResources(ctx, filePath, namespace)` | Applies dependent resources (e.g., for ref-objects) |
| `updateAppNamespaceReferences(app, namespace)` | Updates namespace references inside Application |

#### Constants

```go
const (
    AppRunningTimeout = 5 * time.Minute   // Timeout for app to become running
    PollInterval      = 5 * time.Second   // Polling interval for status checks
)
```

### Test Data

Test data is organized in `test/builtin-definition-example/` with YAML files for each definition type.

#### Example: Component Test Application

```yaml
# test/builtin-definition-example/components/webservice.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: website
spec:
  components:
    - name: frontend
      type: webservice
      properties:
        image: oamdev/testapp:v1
        cmd: ["node", "server.js"]
        ports:
          - port: 8080
            expose: true
```

#### Example: Trait Test Application

```yaml
# test/builtin-definition-example/trait/sidecar.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: vela-app-with-sidecar
spec:
  components:
    - name: log-gen-worker
      type: worker
      properties:
        image: busybox
        cmd: ["/bin/sh", "-c", "while true; do echo hello; sleep 1; done"]
      traits:
        - type: sidecar
          properties:
            name: count-log
            image: busybox
            cmd: ["/bin/sh", "-c", "tail -f /var/log/date.log"]
```

#### Example: Policy Test Application

```yaml
# test/builtin-definition-example/policies/topology.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: basic-topology
spec:
  components:
    - name: nginx-basic
      type: webservice
      properties:
        image: nginx
  policies:
    - name: topology-hangzhou-clusters
      type: topology
      properties:
        clusters: ["local"]
```

#### Example: Workflow Step Test Application

```yaml
# test/builtin-definition-example/workflowsteps/suspend.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: first-vela-workflow
spec:
  components:
    - name: express-server
      type: webservice
      properties:
        image: oamdev/hello-world
        port: 8000
  workflow:
    steps:
      - name: manual-approval
        type: suspend
        properties:
          duration: "30s"
      - name: express-server
        type: apply-component
        properties:
          component: express-server
```

---

## CI/CD Integration

### GitHub Actions Workflows

Four separate workflows run tests for each definition type:

| Workflow | File | Trigger |
|----------|------|---------|
| Test Component Definitions | `test-component-definitions.yaml` | Push/PR to `go.mod`, `go.sum` |
| Test Trait Definitions | `test-trait-definitions.yaml` | Push/PR to `go.mod`, `go.sum` |
| Test Policy Definitions | `test-policy-definitions.yaml` | Push/PR to `go.mod`, `go.sum` |
| Test Workflow Step Definitions | `test-workflowstep-definitions.yaml` | Push/PR to `go.mod`, `go.sum` |


### Reusable Setup Action

The `.github/actions/setup-vela-environment/action.yaml` composite action handles all setup:

#### Setup Steps

1. **Set up Go** - Installs Go 1.23
2. **Set up Kubernetes (Kind)** - Creates a Kind cluster named `vela-test`
3. **Download and install Vela CLI** - Installs latest release for initial setup
4. **Install KubeVela** - Runs `vela install` to deploy KubeVela to the cluster
5. **Extract KubeVela repository info from go.mod** - Parses replace directive to find fork/commit
6. **Clone KubeVela repository and checkout commit** - Clones the specified fork at the exact commit
7. **Build Vela CLI from source** - Runs `make vela-cli` to build the CLI
8. **Uninstall built-in definitions** - Removes default definitions to test module definitions
9. **Install definitions from defkit** - Runs `vela def apply-module .` to install this module's definitions
10. **Install Ginkgo** - Installs Ginkgo CLI for running tests

#### Dynamic KubeVela Version

The action extracts the KubeVela version from `go.mod`:

```bash
# Example go.mod replace directive:
replace github.com/oam-dev/kubevela => github.com/kubevela/kubevela v0.0.0-20250115123456-abc123def456

# Extracted:
# - Repository: kubevela/kubevela
# - Commit SHA: abc123def456
```

This ensures tests always run against the exact KubeVela version the module depends on.

---

## Running Tests Locally

### Prerequisites

1. **Go 1.23+** installed
2. **kubectl** configured with a Kubernetes cluster
3. **KubeVela** installed on the cluster
4. **Ginkgo CLI** installed

### Install Ginkgo

```bash
make install-ginkgo
# or
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Run All Tests

```bash
make test-e2e
```

### Run Specific Test Categories

```bash
# Component definitions only
make test-e2e-components

# Trait definitions only
make test-e2e-traits

# Policy definitions only
make test-e2e-policies

# Workflow step definitions only
make test-e2e-workflowsteps
```

### Configure Parallelism

```bash
# Run with 8 parallel processes (default: 4)
make test-e2e-components PROCS=8

# Run serially (useful for debugging)
make test-e2e-components PROCS=1
```

### Configure Timeout

```bash
# Set 60 minute timeout (default: 30m)
make test-e2e-components E2E_TIMEOUT=60m
```

### Custom Test Data Path

```bash
# Use custom test data directory
make test-e2e-components TESTDATA_PATH=/path/to/custom/test-data
```

### Run with Ginkgo Directly

```bash
# Run with verbose output
ginkgo -v --label-filter="components" ./test/e2e/...

# Run specific test by name
ginkgo -v --focus="webservice" ./test/e2e/...

# Run with race detector
ginkgo -v --race --label-filter="traits" ./test/e2e/...
```

---

## How It Works

### Test Execution Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                     For each YAML test file:                     │
├─────────────────────────────────────────────────────────────────┤
│  1. Read Application from YAML file                              │
│     └─► readAppFromFile(filePath)                                │
│                                                                   │
│  2. Generate unique namespace                                     │
│     └─► uniqueNamespaceForApp(app.Name)                          │
│     └─► Example: "e2e-website" for app "website"                 │
│                                                                   │
│  3. Create namespace                                              │
│     └─► ensureNamespace(ctx, uniqueNs)                           │
│                                                                   │
│  4. Clean up any existing application                             │
│     └─► cleanupExistingApplication(ctx, app)                     │
│                                                                   │
│  5. Apply prerequisite resources (if any)                         │
│     └─► applyPrerequisitesIfAny(ctx, filePath, uniqueNs)         │
│     └─► Example: Deployment/Service for ref-objects tests        │
│                                                                   │
│  6. Apply the Application                                         │
│     └─► k8sClient.Create(ctx, app)                               │
│                                                                   │
│  7. Wait for Application to reach "running" status                │
│     └─► waitForApplicationRunning(ctx, appName, uniqueNs)        │
│     └─► Polls every 5s, timeout 5 minutes                        │
│                                                                   │
│  8. Delete namespace (cleanup)                                    │
│     └─► deleteNamespace(ctx, uniqueNs)                           │
│                                                                   │
│  ✅ Test passes if app reaches "running" state                    │
│  ❌ Test fails if timeout or error status                         │
└─────────────────────────────────────────────────────────────────┘
```

### Parallel Execution

Ginkgo runs tests in parallel processes. Each test:

1. Creates its own unique namespace (e.g., `e2e-website`, `e2e-worker`)
2. Deploys applications in isolation
3. Cleans up its namespace after completion

This allows multiple tests to run simultaneously without interference.

### Prerequisite Resources

Some tests require existing resources (e.g., `ref-objects` references existing Deployments):

```yaml
# ref-objects.yaml contains multiple documents:
---
apiVersion: apps/v1
kind: Deployment        # Prerequisite resource
metadata:
  name: ref-deployment
spec:
  # ...
---
apiVersion: core.oam.dev/v1beta1
kind: Application       # Test application that references the Deployment
metadata:
  name: ref-objects-test
spec:
  components:
    - name: ref-comp
      type: ref-objects
      properties:
        objects:
          - name: ref-deployment
            kind: Deployment
```

The framework automatically:
1. Detects multi-document YAMLs with non-Application resources
2. Applies prerequisite resources first
3. Waits briefly for them to settle
4. Then applies the Application

---

## Adding New Tests

### 1. Create a Test Application YAML

Create a YAML file in the appropriate directory:

```yaml
# test/builtin-definition-example/components/my-component.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: my-component-test
spec:
  components:
    - name: my-app
      type: my-component
      properties:
        image: nginx
        # Add required properties for your component
```

### 2. Test is Automatically Picked Up

The E2E framework automatically discovers all `.yaml` files in the test data directories. No code changes needed!

### 3. Run the Test

```bash
# Run all component tests (includes your new test)
make test-e2e-components

# Or run just your test
ginkgo -v --focus="my-component" ./test/e2e/...
```

---


