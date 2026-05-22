---
name: "program-doanot-skills"
description: "Defines universal coding principles for AI agents. Invoke at the start of any programming task and keep applying it through implementation, debugging, and review."
---

# Program Doanot Skills

Use this skill as a default operating mode for all coding tasks.

## When to invoke

Invoke this skill whenever:
- A user asks for any code change, bug fix, refactor, or feature work
- You are about to make assumptions in unclear requirements
- You are planning implementation steps
- You are about to claim completion and need verification

## Core principles
### 1. Preserve Existing Behavior
- When developing new code, do not break existing functionality unless explicitly requested
- Keep existing business logic correct and running as expected unless explicitly requested to change it
- Only modify or replace existing functionality when the requirement explicitly asks for it
- Verify no regression is introduced before completion

### 2. Optimize for Performance and Complexity
- Choose the lowest practical time complexity for the requirement
- Prefer `O(1)` over `O(n)` when both are feasible
- Prefer `O(n)` over `O(n log n)` when both are feasible
- Prefer `O(n log n)` over `O(n^2)` when both are feasible
- Call out any unavoidable complexity tradeoff explicitly

### 3. Avoid Infinite Loops
- Never introduce unbounded loops without clear termination conditions
- Ensure retries, polling, and recursion have explicit stop conditions
- Add guards (max attempts, timeout, or break conditions) where needed

### 4. Make Surgical Changes
- Touch only files and lines required by the request
- Do not refactor unrelated areas
- Preserve existing style and patterns unless asked to change them
- Remove only dead code introduced by your own edits

### 5. Keep It Simple
- Implement the smallest change that satisfies the request
- Avoid speculative abstractions, new layers, and over-generalization
- Prefer straightforward code over "future-proof" complexity

### 6. Fix the Root Cause of Bugs
- If the program has an error or bug, do not try to hide it or “patch around it” to make the symptom disappear
- Identify the real root cause and fix that cause directly
- Only treat it as an external constraint when the root cause is outside your codebase (e.g., third-party tools, databases, or other external components)

### 7. Follow Existing Protocol Conventions (No Duplicates)
- Protocol definitions must follow the existing protocol/spec conventions in the codebase
- If a structure/message/type already exists somewhere, do not redefine it in another place
- Reuse the same shared type to keep a single source of truth and avoid split-brain updates

### 8. Observability and Operability (Log Key Steps)
- Log at important steps and on errors (e.g., Redis errors, database errors, upstream/API call failures)
- Logs must include key identifiers where applicable (e.g., `key`, `user_id`, `request_id`) to support debugging and operations
- Use structured logging when available; include enough context to locate the failing dependency and operation
- Never log secrets or sensitive payloads (tokens, passwords, full PII)

### 9. DRY Principle (Do Not Repeat Yourself)
- Do not write duplicate code. If a piece of code or feature is used in multiple places, extract it into a shared function to ensure modularity and encapsulation.
- Intra-module: If code is used in 2 or more places within the same module, extract it into a common function call.
- Cross-module: If code has the same functionality in 2 or more modules (e.g., generating JWT tokens, database access, message queue pub/sub, Redis access), extract it into a shared module or common code for reuse across modules.
