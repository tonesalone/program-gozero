# CLAUDE Integration Guide

This guide explains how to use `program-doanot-skills` in Claude-oriented workflows.

## Purpose

Use this skill as a default coding policy to reduce regressions, avoid unnecessary complexity, and enforce safe loop behavior.

## Recommended Usage Pattern

Before implementation starts, include an instruction such as:

```text
Apply program-doanot-skills as baseline constraints for this coding task.
```

During execution, ensure the agent:

- Preserves existing functionality and business logic
- Chooses the lowest practical complexity
- Avoids unbounded loops and adds termination guards
- Makes surgical changes only
- Keeps implementation minimal and simple
- Reuses existing protocol/type definitions; do not duplicate structures
- Logs key steps and errors with identifying context (e.g., `key`, `user_id`, `request_id`)

## Suggested Prompt Snippet

```text
Use program-doanot-skills for this task.
Non-negotiable constraints:
1) Do not break existing behavior.
2) Prefer lower complexity solutions when feasible.
3) No infinite loops; retries/polling/recursion require termination guards.
4) Keep edits minimal and limited to relevant files.
5) Avoid speculative abstractions.
6) Follow existing protocol conventions; never redefine a type that already exists.
7) Add actionable logs at key steps and on failures (Redis/DB/API), including key identifiers.
```

## Verification Expectations

Before claiming completion, confirm:

- Existing behavior remains correct (no regression observed)
- Complexity choice is justified
- Loop-related logic has explicit stop conditions
- Changed lines are directly traceable to the request
- Protocol/type definitions are reused and not duplicated
- Key failure points are observable via logs with identifiers
