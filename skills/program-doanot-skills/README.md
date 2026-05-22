# Program Doanot Skills

A universal coding-discipline skill for AI coding agents.

## Overview

`program-doanot-skills` defines a strict baseline for safe, efficient, and regression-resistant coding.  
Use it as a default rule set for feature work, bug fixes, refactors, and implementation planning.

## Features

- Protects existing behavior and business logic from unintended breakage
- Enforces practical performance and complexity choices
- Prevents infinite loops with explicit termination requirements
- Promotes focused, minimal, request-driven code changes
- Encourages simple implementations over speculative abstractions
- Reuses existing protocol/type definitions to keep a single source of truth
- Improves observability with actionable logs at key failure points
- Enforces DRY principle by extracting duplicate code into shared functions and common modules

## When To Invoke

Invoke this skill when:

- Starting any programming task
- Modifying existing code paths
- Introducing loops, retries, polling, or recursive logic
- Making tradeoffs between multiple algorithmic approaches
- Preparing final validation before declaring completion

## Core Rules

1. Preserve Existing Behavior
- New code must not break existing functionality.
- Existing business logic must remain correct.
- Verify no regression before completion.

2. Optimize for Performance and Complexity
- Choose the lowest practical time complexity.
- Prefer `O(n)` over `O(n log n)` when feasible.
- Prefer `O(n log n)` over `O(n^2)` when feasible.
- State unavoidable complexity tradeoffs explicitly.

3. Avoid Infinite Loops
- Never introduce unbounded loops without clear exit conditions.
- Retries, polling, and recursion must include stop conditions.
- Add safeguards such as max attempts, timeout, or break logic.

4. Make Surgical Changes
- Edit only files and lines required by the task.
- Do not refactor unrelated areas.
- Preserve local style and patterns unless asked to change them.
- Remove only dead code introduced by your own edits.

5. Keep It Simple
- Implement the smallest change that meets the requirement.
- Avoid speculative abstractions and unnecessary layers.
- Prefer straightforward code over "future-proof" complexity.

6. Follow Existing Protocol Conventions (No Duplicates)
- Protocol definitions must follow the existing protocol/spec conventions in the codebase.
- If a structure/message/type already exists somewhere, do not redefine it in another place.
- Reuse the same shared type to keep a single source of truth.

7. Observability and Operability (Log Key Steps)
- Log at important steps and on errors (e.g., Redis errors, database errors, upstream/API call failures).
- Logs must include key identifiers where applicable (e.g., `key`, `user_id`, `request_id`).
- Never log secrets or sensitive payloads (tokens, passwords, full PII).

8. DRY Principle (Do Not Repeat Yourself)
- Do not write duplicate code.
- If a feature or logic is used in multiple places, extract it into a shared function for modularity and encapsulation.
- Intra-module: Extract common function calls for code used in 2+ places within a module.
- Cross-module: Extract shared code (e.g., JWT generation, DB access, MQ pub/sub, Redis access) used in 2+ modules into a common module.

## Usage Example

Example invocation in agent instruction:

```text
Use program-doanot-skills for this task.
Before coding, confirm no regression risk and choose the lowest practical complexity.
Any loop/retry logic must include explicit termination guards.
```

## Quick Checklist

- Does this change preserve existing behavior?
- Is the chosen approach the lowest practical complexity?
- Can any loop/retry/polling path terminate safely?
- Are edits limited to required files and lines?
- Is the implementation minimal and non-speculative?

## Limitations

- This skill provides behavioral guardrails, not framework-specific APIs.
- It does not replace domain requirements, architecture specs, or test strategy.
- It should be combined with project-level conventions when available.
