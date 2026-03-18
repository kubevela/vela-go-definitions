# E2E Test Automation Framework

## Overview

The E2E test framework validates that all 77 KubeVela X-Definitions (components, traits, policies, workflow steps) work correctly when deployed to a live Kubernetes cluster. It uses a two-layer validation approach: auto-derived checks for every definition, plus optional `.expect.yaml` files for definition-specific assertions.

### Key Features

- **77 definitions tested** across 4 types (8 components, 29 traits, 9 policies, 31 workflow steps)
- **Two-layer validation**: auto-derived baseline + optional extra expectations
- **Multi-app support**: tests with multiple Applications (e.g., shared-resource, depends-on-app)
- **Parallel execution**: Ginkgo multi-process with isolated namespaces
- **One-command setup**: `make e2e-setup` creates a k3d cluster with everything installed
- **Failure diagnostics**: workflow status, kubectl describe, pod logs on failure

---

## Quick Start

```bash
# Set up local test environment (k3d + KubeVela + definitions)
make e2e-setup

# Run all tests
make test-e2e

# Run by category
make test-e2e-components
make test-e2e-traits
make test-e2e-policies
make test-e2e-workflowsteps

# Tear down
make e2e-teardown
```

Prerequisites: docker, k3d, kubectl, vela CLI, Go 1.23+

---

## Directory Structure

```
test/
  e2e/
    e2e_suite_test.go          # Ginkgo suite bootstrap
    definition_e2e_test.go     # Table-driven test generator for all 4 types
    helpers_test.go            # All test logic: runner, auto-validate, expectations
  builtin-definition-example/
    applications/              # Test inputs (Application YAMLs)
      components/              # 8 component tests
      trait/                   # 29 trait tests
      policies/                # 9 policy tests
      workflowsteps/           # 31 workflow step tests
    expectations/              # Extra validation (optional, additive)
      trait/                   # Trait-specific checks
      policies/                # Policy-specific checks
      workflowsteps/           # Workflow step output checks
```

---

## Architecture

### Test Runner (`definition_e2e_test.go`)

A single file generates all test suites from a table-driven config:

```go
var suites = []definitionTestSuite{
    {label: "components",    subdir: "applications/components",    descName: "Component"},
    {label: "traits",        subdir: "applications/trait",         descName: "Trait"},
    {label: "policies",      subdir: "applications/policies",      descName: "Policy"},
    {label: "workflowsteps", subdir: "applications/workflowsteps", descName: "WorkflowStep", skipTests: skipWorkflowStepTests},
}
```

Each YAML file in the test data directory becomes a Ginkgo `It` block. Tests are auto-discovered — no code changes needed to add new tests.

### Test Execution Flow

For each test YAML file:

1. **Parse** all Applications from the file (supports multi-doc YAML with multiple apps)
2. **Create** isolated namespace (`e2e-{appname}`)
3. **Apply prerequisites** — non-Application resources (Deployments, Services, ConfigMaps) with polling for readiness
4. **Apply all Applications** sequentially, waiting for each to reach `phase: running`
5. **Auto-validate** (Layer 1):
   - All workflow steps have `phase: succeeded`
   - Component resources exist (Deployment, DaemonSet, StatefulSet, Job, CronJob) with correct image
6. **Extra validation** (Layer 2) from `.expect.yaml` if it exists:
   - Resource field assertions (dot-path with array indexing and bracket key notation)
   - Workflow step message assertions
7. **Cleanup** — delete apps (clears finalizers), then delete namespace
8. **Diagnostics** on failure — app status, workflow steps, vela status, kubectl describe, pod logs

### Two-Layer Validation

**Layer 1 — Auto-derived (all 77 definitions, zero config):**

Every test automatically validates:
- Workflow steps all reached `phase: succeeded` (including sub-steps)
- Component resources exist based on type mapping:

| Component Type | K8s Resource | Image Path |
|---------------|-------------|------------|
| `webservice`, `worker` | Deployment | `spec.template.spec.containers[0].image` |
| `daemon` | DaemonSet | `spec.template.spec.containers[0].image` |
| `statefulset` | StatefulSet | `spec.template.spec.containers[0].image` |
| `task` | Job | `spec.template.spec.containers[0].image` |
| `cron-task` | CronJob | `spec.jobTemplate.spec.template.spec.containers[0].image` |
| `k8s-objects`, `ref-objects` | (varied, skipped) | — |

Resource names are resolved from `status.appliedResources` (not guessed).

**Layer 2 — `.expect.yaml` files (25 files for definition-specific checks):**

Only needed when testing something the auto-derive can't check. Format:

```yaml
# Validate K8s resource fields
expectations:
  - apiVersion: apps/v1
    kind: Deployment
    name: nginx-app
    fields:
      spec.replicas: 3
      spec.template.metadata.annotations["prometheus.io/scrape"]: "true"
      spec.template.spec.containers[0].env[0].name: "LOG_LEVEL"

# Validate workflow step status
workflowSteps:
  - name: message
    phase: succeeded
    messageContains: "All addons have been enabled"
```

Supported path syntax:
- Dot notation: `spec.template.spec.containers`
- Array indexing: `containers[0].image`
- Bracket keys (for dots/slashes in names): `annotations["app.example.com/owner"]`

### Multi-App Support

Some tests contain multiple Applications in a single YAML file (e.g., `shared-resource.yaml`, `depends-on-app.yaml`). The framework:

1. Parses ALL Applications from the file
2. Applies them sequentially, waiting for each to reach running
3. Auto-validates the **last** application (the "main" one; earlier apps are dependencies)

### Skipped Tests

Tests requiring external infrastructure are skipped:

| Test | Reason |
|------|--------|
| `deploy-cloud-resource.yaml` | Requires Alibaba RDS + multi-cluster |
| `share-cloud-resource.yaml` | Requires Alibaba RDS + multi-cluster |
| `generate-jdbc-connection.yaml` | Requires Alibaba RDS |
| `apply-terraform-config.yaml` | Requires Terraform provider credentials |
| `apply-terraform-provider.yaml` | Requires Terraform provider credentials |
| `build-push-image.yaml` | Requires external container registry (ttl.sh) |
| `check-metrics.yaml` | Requires external Prometheus endpoint |
| `restart-workflow.yaml` | Self-restarting workflow, can't validate with single-shot framework |

---

## CI/CD Integration

### GitHub Actions Workflow (`test-definitions.yaml`)

A single workflow with 4 parallel jobs (one per definition type):

| Job | Label Filter | What it tests |
|-----|-------------|---------------|
| `test-components` | `components` | 8 component definitions |
| `test-traits` | `traits` | 29 trait definitions |
| `test-policies` | `policies` | 9 policy definitions |
| `test-workflowsteps` | `workflowsteps` | 31 workflow step definitions |

### Setup Action (`.github/actions/setup-vela-environment`)

Reusable composite action that:

1. Sets up Go + k3d cluster
2. Downloads and installs vela CLI (latest release)
3. Installs KubeVela on the cluster
4. Extracts kubevela fork/commit from `go.mod` replace directive
5. Clones and builds vela CLI from source (for `apply-module` support)
6. Uninstalls built-in CUE definitions
7. Installs defkit definitions via `vela def apply-module .`
8. Installs Ginkgo

The built-from-source CLI uses `cmd/register/main.go` (fast registry path) to discover all 77 definitions.

---

## Local Development

### Make Targets

| Target | Description |
|--------|-------------|
| `e2e-setup` | Create k3d cluster, install KubeVela, install definitions, install Ginkgo |
| `e2e-teardown` | Delete the k3d cluster |
| `test-e2e` | Run all E2E tests (components + traits + policies + workflowsteps) |
| `test-e2e-components` | Run component tests only |
| `test-e2e-traits` | Run trait tests only |
| `test-e2e-policies` | Run policy tests only |
| `test-e2e-workflowsteps` | Run workflow step tests only |
| `cleanup-e2e-namespaces` | Delete all `e2e-*` namespaces |
| `force-cleanup-e2e-namespaces` | Force-delete stuck terminating namespaces |

Each `test-e2e-*` target automatically runs `force-cleanup-e2e-namespaces` first.

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PROCS` | 10 | Parallel Ginkgo processes |
| `E2E_TIMEOUT` | 10m | Total test suite timeout |
| `TESTDATA_PATH` | `test/builtin-definition-example` | Test data directory |
| `E2E_CLUSTER` | `e2e-test` | k3d cluster name |

### Running Individual Tests

```bash
# Run a single test by name
TESTDATA_PATH=test/builtin-definition-example \
  ginkgo -v --timeout=5m --focus="webservice.yaml" --label-filter="components" ./test/e2e/...

# Run serially for debugging
make test-e2e-components PROCS=1
```

---

## Adding New Tests

### 1. Add Application YAML

Create a YAML file in the appropriate `applications/` subdirectory:

```yaml
# test/builtin-definition-example/applications/trait/my-trait.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: my-trait-example
spec:
  components:
    - name: nginx-app
      type: webservice
      properties:
        image: nginx:latest
      traits:
        - type: my-trait
          properties:
            key: value
```

The test is auto-discovered — no code changes needed. It will automatically:
- Verify app reaches running
- Verify all workflow steps succeeded
- Verify the Deployment exists with correct image

### 2. Add Extra Expectations (Optional)

If the trait/policy/step has specific effects to validate, create an `.expect.yaml` file:

```yaml
# test/builtin-definition-example/expectations/trait/my-trait.expect.yaml
expectations:
  - apiVersion: apps/v1
    kind: Deployment
    name: nginx-app
    fields:
      spec.template.metadata.labels.my-key: "value"
```

### 3. Run

```bash
ginkgo -v --timeout=5m --focus="my-trait.yaml" --label-filter="traits" ./test/e2e/...
```

---

## Troubleshooting

### Stuck Namespaces

```bash
make force-cleanup-e2e-namespaces
```

### inotify Errors on macOS (k3d nodes not starting)

The `e2e-setup` target fixes this automatically, but if needed manually:

```bash
docker run --rm --privileged alpine:latest sh -c \
  "sysctl -w fs.inotify.max_user_watches=524288 && sysctl -w fs.inotify.max_user_instances=512"
```

### Test Failures — Reading Diagnostics

On failure, the framework prints:
- Application phase and workflow step details
- `vela status` output
- `kubectl describe app` output
- Pod listing in the namespace

### Definitions Not Installing

If `vela def apply-module .` fails, the fallback is:

```bash
make generate
for f in vela-templates/definitions/*/*.cue; do vela def apply "$f"; done
```
