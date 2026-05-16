# proto-to-gozero-api

Generate go-zero `.api` files from protobuf `.proto` definitions.

## What it does

- Converts protobuf `message` to go-zero `type`
- Converts protobuf `rpc` to go-zero API endpoints
- Applies method/route mapping from annotations or agreed conventions
- Preserves compatibility expectations unless change is explicitly requested

## When to use

- Creating initial go-zero API spec from existing proto contracts
- Regenerating `.api` after proto service/message updates
- Standardizing proto-to-api conversion behavior across a team

## Files

- `SKILL.md`: English skill definition
- `SKILL.zh.md`: Chinese skill definition
- `CLAUDE.md`: Claude-oriented usage guidance
- `CURSOR.md`: Cursor-oriented usage guidance

## Quick usage

Ask the agent to:

```text
Use proto-to-gozero-api.
Generate <target>.api from <source>.proto.
Prefer proto HTTP annotations; if missing, ask for fallback route policy before generation.
```

