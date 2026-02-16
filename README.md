# vela-definitions

A collection of KubeVela X-Definitions.

## Overview

This module contains Go-based KubeVela X-Definitions that can be applied to any KubeVela cluster.

## Directory Structure

- **components/** - ComponentDefinitions for workload types
- **traits/** - TraitDefinitions for operational behaviors
- **policies/** - PolicyDefinitions for application policies
- **workflowsteps/** - WorkflowStepDefinitions for delivery workflows

## Usage

### Apply all definitions

```bash
vela def apply-module github.com/anoop2811/vela-definitions
```

### List definitions

```bash
vela def list-module github.com/anoop2811/vela-definitions
```

### Validate definitions

```bash
vela def validate-module github.com/anoop2811/vela-definitions
```

### Apply with namespace

```bash
vela def apply-module github.com/anoop2811/vela-definitions --namespace my-namespace
```

### Dry-run (preview without applying)

```bash
vela def apply-module github.com/anoop2811/vela-definitions --dry-run
```

## Adding New Definitions

1. Create a new Go file in the appropriate directory
2. Add an init() function that registers your definition
3. Use the defkit package fluent API to define your component/trait/policy/workflow-step
4. Run `go mod tidy` to update dependencies
5. Validate with `vela def validate-module .`

Example component definition:

```go
package components

import "github.com/oam-dev/kubevela/pkg/definition/defkit"

func init() {
    defkit.Register(MyComponent())
}

func MyComponent() *defkit.ComponentDefinition {
    image := defkit.String("image").Required().Description("Container image")
    replicas := defkit.Int("replicas").Default(1).Description("Number of replicas")

    return defkit.NewComponent("my-component").
        Description("My custom component").
        Workload("apps/v1", "Deployment").
        Params(image, replicas).
        Template(myComponentTemplate)
}

func myComponentTemplate(tpl *defkit.Template) {
    vela := defkit.VelaCtx()
    image := defkit.String("image")
    replicas := defkit.Int("replicas")

    deployment := defkit.NewResource("apps/v1", "Deployment").
        Set("spec.replicas", replicas).
        Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
        Set("spec.template.spec.containers[0].name", vela.Name()).
        Set("spec.template.spec.containers[0].image", image)

    tpl.Output(deployment)
}
```

## Testing

This module includes a comprehensive E2E test automation framework that validates all definitions against a live KubeVela cluster.

### Quick Start

```bash
# Install Ginkgo CLI
make install-ginkgo

# Run all E2E tests
make test-e2e

# Run specific test categories
make test-e2e-components    # Component definitions
make test-e2e-traits        # Trait definitions
make test-e2e-policies      # Policy definitions
make test-e2e-workflowsteps # Workflow step definitions
```

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PROCS` | 4 | Number of parallel test processes |
| `E2E_TIMEOUT` | 30m | Test timeout duration |
| `TESTDATA_PATH` | `test/builtin-definition-example` | Path to test data |

### CI/CD

GitHub Actions workflows automatically run E2E tests on:
- Push to `go.mod` or `go.sum`
- Pull requests modifying `go.mod` or `go.sum`
- Manual workflow dispatch

For detailed documentation, see [E2E Test Automation Guide](docs/E2E_TEST_AUTOMATION.md).

## Version

Version is automatically derived from git tags. Use semantic versioning tags (e.g., v1.0.0) to set the module version.

## License

Apache License 2.0
