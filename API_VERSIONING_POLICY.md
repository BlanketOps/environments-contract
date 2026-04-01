
# BlanketOps Environments Contract Versioning & Promotion Policy

## Overview

This document defines the versioning and promotion model for all BlanketOps API contracts.

The goal is to ensure:

* Safe evolution of APIs
* Backward compatibility where required
* Clear expectations for consumers
* Strong governance of domain contracts

All APIs follow a staged promotion model:

```
v1alpha1 → v1beta1 → v1
```

---

## Version Stages

### 1. v1alpha1 — Exploration Stage

**Purpose**

* Rapid iteration
* Domain discovery
* Early experimentation

**Rules**

* Breaking changes are allowed
* Fields may be added, removed, or renamed
* Message structure may change
* No backward compatibility guarantees

**Usage**

* Internal development
* Early prototyping
* Experimental features

**Principle**

> The shape of the domain is still being discovered

---

### 2. v1beta1 — Stabilization Stage

**Purpose**

* Stabilize structure
* Validate real-world usage
* Prepare for production readiness

**Rules**

* Breaking changes are strongly discouraged
* Fields MUST NOT be removed
* Changes should be additive
* Backward compatibility is expected

**Allowed Changes**

* Add optional fields
* Extend enums (append-only)
* Clarify semantics (without breaking meaning)

**Not Allowed**

* Removing fields
* Renaming fields
* Changing field types

**Usage**

* Internal production usage
* Staging environments
* Wider adoption

**Principle**

> The shape is known, now we validate behavior

---

### 3. v1 — Stable Contract

**Purpose**

* Production-grade API
* External consumption
* Long-term stability

**Rules**

* No breaking changes allowed
* Strict backward compatibility
* Only additive changes permitted
* Deprecated fields must remain

**Allowed Changes**

* Add new optional fields
* Append enum values

**Not Allowed**

* Removing fields
* Renaming fields
* Changing semantics of existing fields

**Usage**

* Production systems
* External SDKs
* Public APIs

**Principle**

> This contract is trusted and must not break

---

## Promotion Criteria

### Alpha → Beta

Promote when:

* Core structure is stable
* Domain model is understood
* End-to-end workflows function
* No frequent structural changes required

Signal:

> We understand the shape of this resource

---

### Beta → v1

Promote when:

* No breaking changes required over time
* Used in real environments successfully
* Edge cases handled
* Behavior is predictable and stable

Signal:

> This contract is safe for production

---

## Contract Rules (Enforced)

### 1. Field Numbers Are Immutable

Field numbers must never be reused.

```
string name = 1; // permanently reserved
```

---

### 2. Fields Are Deprecated, Not Removed

```
string old_field = 5 [deprecated = true];
```

---

### 3. Enums Are Append-Only

```
enum Phase {
  PHASE_UNSPECIFIED = 0;
  PHASE_PENDING = 1;
  // new values appended only
}
```

---

### 4. Package Defines Version Boundary

```
package blanketops.environments.v1alpha1;
package blanketops.environments.v1beta1;
package blanketops.environments.v1;
```

Each version is isolated and independently evolvable.

---

## Design Principles

* Contracts evolve through stages, not ad hoc changes
* Stability increases as versions progress
* Backward compatibility is intentional, not accidental
* Domain meaning must remain consistent across versions

---

## Summary

| Stage    | Stability | Breaking Changes | Usage               |
| -------- | --------- | ---------------- | ------------------- |
| v1alpha1 | Low       | Allowed          | Internal / Dev      |
| v1beta1  | Medium    | Discouraged      | Staging / Early     |
| v1       | High      | Forbidden        | Production / Public |

---

## Version: v0.1.7

This version introduces:

* Formal contract promotion policy
* Standardized versioning model across all domains
* Governance rules for future evolution

---

## Final Note

BlanketOps contracts are designed to evolve safely.

This policy ensures:

> innovation without instability
>
