# CURSOR Guide

Apply `proto-to-gozero-api` in Cursor sessions when transforming `.proto` definitions to go-zero `.api`.

## Recommended instruction

```text
Use proto-to-gozero-api for this task.
Generate API definitions from the provided proto file.
If method/route mapping is not explicit, ask for the mapping policy first.
Keep compatibility unless the requirement explicitly asks for breaking changes.
```

## Checklist

- Service and RPC parsing is complete
- Message-to-type mapping is complete
- Method + route exists for every endpoint
- Compatibility impact is documented
