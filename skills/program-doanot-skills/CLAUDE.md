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

## Suggested Prompt Snippet

```text
Use program-doanot-skills for this task.
Non-negotiable constraints:
1) Do not break existing behavior.
2) Prefer lower complexity solutions when feasible.
3) No infinite loops; retries/polling/recursion require termination guards.
4) Keep edits minimal and limited to relevant files.
5) Avoid speculative abstractions.
```

## Verification Expectations

Before claiming completion, confirm:

- Existing behavior remains correct (no regression observed)
- Complexity choice is justified
- Loop-related logic has explicit stop conditions
- Changed lines are directly traceable to the request
