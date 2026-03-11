# Roadmap: defkit API Documentation

## Overview

Four phases build a single-page HTML reference for the defkit Go API. Each phase produces a distinct, appendable chunk of docs/index.html. Phase 1 establishes the page shell and navigation. Phase 2 documents definition builders and parameter API. Phase 3 documents template, resource, value, and context APIs. Phase 4 delivers full worked examples per definition type.

## Phases

- [ ] **Phase 1: Page Shell** - HTML structure, CSS, sidebar, navigation skeleton
- [ ] **Phase 2: Definition + Parameter API** - NewComponent/Trait/Policy/WorkflowStep + all param types and chain methods
- [ ] **Phase 3: Template + Resource + Value + VelaCtx** - tpl.Output, resource builders, value expressions, VelaCtx
- [ ] **Phase 4: Full Examples** - One complete worked example per definition type

## Phase Details

### Phase 1: Page Shell
**Goal**: A navigable, styled HTML page skeleton a developer can open in a browser
**Depends on**: Nothing (first phase)
**Requirements**: STRUCT-01, STRUCT-02, STRUCT-03, STRUCT-04, STRUCT-05, FMT-01, FMT-02, FMT-03
**Success Criteria** (what must be TRUE):
  1. Developer opens docs/index.html and sees a dark sidebar with all section links
  2. Clicking a sidebar link scrolls to the correct content section
  3. Code blocks render with syntax highlighting (Go and CUE visually distinct)
  4. Page is usable on a standard desktop viewport without horizontal scrolling
**Plans**: 1 plan

Plans:
- [ ] 01-01-PLAN.md — HTML shell: layout, sidebar, Prism.js syntax highlighting, 8 placeholder sections

### Phase 2: Definition + Parameter API
**Goal**: All four definition builder constructors and every parameter type are documented with Go and CUE examples
**Depends on**: Phase 1
**Requirements**: DEF-01, DEF-02, DEF-03, DEF-04, DEF-05, DEF-06, DEF-07, DEF-08, PARAM-01, PARAM-02, PARAM-03, PARAM-04
**Success Criteria** (what must be TRUE):
  1. Developer can find NewComponent, NewTrait, NewPolicy, NewWorkflowStep with signature and example
  2. AutodetectWorkload() and RawCUE() are documented with when/why to use them
  3. Every scalar, collection, and complex parameter type has a documented entry
  4. Every chain method (Required, Optional, Default, Description, Values, Min, Max, WithFields, Of) shows Go snippet and CUE equivalent side by side
  5. Definition-level chain methods (Description, Workload, Params, Template, AppliesTo, WithImports) are documented per definition type
**Plans**: TBD

### Phase 3: Template + Resource + Value + VelaCtx + Health
**Goal**: The full runtime-construction API is documented — templates, resource builders, value expressions, context accessors, and health/status DSL
**Depends on**: Phase 2
**Requirements**: TPL-01, TPL-02, TPL-03, TPL-04, TPL-05, TPL-06, RES-01, RES-02, RES-03, RES-04, RES-05, RES-06, RES-07, RES-08, VAL-01, VAL-02, VAL-03, VAL-04, VAL-05, VAL-06, CTX-01, HEALTH-01, HEALTH-02, HEALTH-03
**Success Criteria** (what must be TRUE):
  1. Developer can find tpl.Output, tpl.Outputs, tpl.Patch, tpl.UsePatchContainer with examples
  2. tpl.SetRawHeaderBlock, tpl.SetRawOutputsBlock, and tpl.Helper() builder are documented
  3. Every resource builder method (Set, SetIf, ForEach, ForEachWith, ItemBuilder, If/EndIf, Directive, VersionIf) has a Go snippet and CUE equivalent
  4. Value constructors (Lit, Reference, Interpolation) and chain methods are documented including PathExists, ParamRef, ParamNotSet
  5. Let bindings (Let, LetVariable) are documented
  6. Logical operators, string operations, and comparison operators documented
  7. VelaCtx accessor methods documented
  8. Health DSL (Health(), HealthyWhen, preset builders) and Status DSL documented
**Plans**: TBD

### Phase 4: Full Examples
**Goal**: A developer can read a complete end-to-end example for each definition type and immediately replicate the pattern
**Depends on**: Phase 3
**Requirements**: EX-01, EX-02, EX-03, EX-04
**Success Criteria** (what must be TRUE):
  1. A complete Component definition example (webservice-like) is shown with full Go code and init() registration
  2. A complete Trait definition example (env-like patch trait) is shown with PatchContainer usage
  3. A complete Policy definition example is shown with Template method
  4. A complete WorkflowStep definition example is shown end-to-end
**Plans**: TBD

## Progress

**Execution Order:** 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Page Shell | 0/1 | Not started | - |
| 2. Definition + Parameter API | 0/TBD | Not started | - |
| 3. Template + Resource + Value + VelaCtx | 0/TBD | Not started | - |
| 4. Full Examples | 0/TBD | Not started | - |
