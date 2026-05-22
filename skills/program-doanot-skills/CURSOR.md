# Cursor Integration Guide

This document describes how to apply `program-doanot-skills` in Cursor-based coding sessions.

## Goal

Keep coding output safe, efficient, and minimal by applying a fixed set of engineering constraints.

## How To Use In Cursor

Add a task-level instruction that explicitly activates this skill behavior:

```text
Use program-doanot-skills as the baseline coding rules for this task.
```

## Operational Rules

The agent should consistently enforce:

- Regression protection for existing functionality and business logic
- Practical complexity optimization (`O(n)` preferred over `O(n log n)` where feasible)
- Explicit termination safeguards for loops, retries, polling, and recursion
- Minimal, task-scoped, surgical edits only
- Simplicity-first implementation choices
- Follow existing protocol conventions; reuse existing types and avoid duplicate definitions
- Ensure observability: log key steps and failures with identifiers (e.g., `key`, `user_id`, `request_id`)
- Enforce DRY principle: extract duplicate code into shared functions (intra-module) or common code (cross-module for things like JWT, DB, MQ, Redis)

## Recommended Session Checklist

- Confirm expected behavior before changing code
- Define complexity target for the solution path
- Add termination conditions for any repeated flow
- Validate modified areas and check for regressions
- Ensure no unrelated refactor is introduced
- Confirm protocol/type definitions are not duplicated across modules
- Confirm key error paths (Redis/DB/API) emit actionable logs with identifying context
- Confirm no duplicated code is written; extract shared functions and common code when used in multiple places

## Notes

- This guide is skill-level governance, not a replacement for project architecture docs.
- For best results, combine with repository-specific standards and tests.
