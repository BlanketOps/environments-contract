# Blanketops Environments Contract

## Overview.

`Blanketops Environments Contract` defines the canonical contracts for BlanketOps Environments.

This repository contains `versioned protobuf schemas` that describe:

- environment intent.
- build and deployment workflows.
- routing and service units.
- external events and sources.

These contracts are the source of truth for how BlanketOps components communicate and reason about environments.

---

## Repository Structure

```tree
blanketops/
├── environments/
│ ├── v1alpha1/
│ ├── v1beta1/
│ └── v1/
├── events/
│ ├── v1alpha1/
│ ├── v1beta1/
│ └── v1/
└── networks/
│ ├── v1alpha1/
│ ├── v1beta1/
│ ├── v1/
│ 
└── sources/
  ├── v1alpha1/
  ├── v1beta1/
  └── v1/
```

Each API group:

- is versioned explicitly.
- contains .proto files as the contract.
- includes generated .pb.go files for Go consumers.
- Version directories are append-only once stabilized.

## Versioning Model

This repository follows semantic versioning at the contract level.

### API versions

- `v1alpha1` — experimental, breaking changes allowed
- `v1beta1` — feature complete, limited breaking changes
- `v1` — stable, backward-compatible

Once an API reaches `v1`:

- fields are not removed.
- semantics are preserved.
- evolution happens via additive changes only.

Breaking changes require a `new API version.`

## Code Generation

This repository uses Buf for protobuf generation and consistency.

Relevant files:

- `buf.gen.yaml` — generation configuration
- `buf.gen` — generation output settings

Generated Go code is committed intentionally to:

- ensure reproducible builds.
- simplify downstream consumption.
- avoid forcing consumers to run codegen.

## Tooling

Common tasks are exposed via the Makefile.

Typical workflows include:

- validating protobuf definitions.
- regenerating Go bindings.
- ensuring schema consistency.

Refer to:

```Magefile
mage help
```

for available targets.

### Consumers

This repository is consumed by:

- BlanketOps operators and controllers.
- Internal policy and resolution engines.
- Tooling that requires a stable environment contract.
