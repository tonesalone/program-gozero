# CLAUDE Guide

Use `proto-to-gozero-api` when converting protobuf contracts into go-zero `.api` specs.

## Recommended prompt

```text
Use proto-to-gozero-api.
Read <source>.proto and generate <target>.api.
Use proto HTTP annotations first.
If mapping is ambiguous, stop and ask before generating.
Do not introduce breaking route or type changes unless explicitly requested.
```

## Expected behavior

- Deterministic mapping from proto to `.api`
- Clear clarification when route/method policy is missing
- Explicit compatibility notes for any breaking impact
