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
