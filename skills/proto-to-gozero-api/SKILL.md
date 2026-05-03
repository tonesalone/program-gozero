---
name: "proto-to-gozero-api"
description: "Generates go-zero .api specs from protobuf .proto files. Invoke when user asks to convert proto RPC/message definitions into go-zero API routes, types, and service blocks."
---

# Proto To go-zero API

Convert protobuf `.proto` definitions into go-zero `.api` files with predictable mapping rules and minimal ambiguity.

## When to invoke

Invoke this skill when:
- User asks to generate a go-zero `.api` file from a `.proto` file
- User wants RPC/service/message definitions mapped to go-zero API schema
- User needs route, method, request, and response definitions derived from protobuf
- Existing `.api` needs to be regenerated after `.proto` changes

## Inputs required

- Source `.proto` file path
- Target `.api` file path
- Routing convention:
  - Prefer explicit HTTP annotations in proto if present
  - Otherwise use agreed fallback mapping (must be stated)
- Naming policy:
  - Keep proto naming as-is, or
  - Convert to project naming style (must be stated)

## Output requirements

The generated `.api` should include:
- `type` blocks mapped from protobuf messages used by service RPCs
- Service/API declaration blocks for RPC endpoints
- HTTP method + route path for each endpoint
- Request/response type references aligned with proto definitions

Strict scope rule:
- Convert only data structures involved by interfaces under the proto `service` block
- Do not convert unrelated standalone messages not referenced by service RPC request/response chains

## Core mapping rules

### 1. Message to type mapping
- Map protobuf `message` to go-zero `type`
- Preserve field order unless project convention requires reordering
- Keep semantic field names; avoid unnecessary renaming
- Type naming follows sample style: `AddReq` -> `addReq`, `CheckResp` -> `checkResp`
- Request types must include both `json` and `form` tags (sample style)
- Response types should include `json` tags at minimum

Example:

Proto:
```proto
message AddReq {
  string book = 1;
  int64 price = 2;
}
```

API:
```text
addReq {
  book string `json:"book" form:"book"`
  price int64 `json:"price" form:"price"`
}
```

### 2. RPC to endpoint mapping
- Map protobuf `rpc` to one API handler endpoint
- Keep request/response type alignment exactly with proto RPC signature
- If route or method is missing, ask for mapping policy before generating
- Handler naming follows RPC name: `Add` -> `@handler AddHandler`

### 3. Route and method strategy
- First choice: use explicit proto HTTP annotation comment if available
- Fallback: apply team/project convention consistently
- Never guess mixed conventions in one file
- Parse annotation style like: `// http POST /cs/v1/modname/add`
- Map `POST` -> `post`, `GET` -> `get` in go-zero API method keyword
- Keep path exactly aligned with annotation path
- Module segment (e.g. `modname`) must align with the proto module/package context

Example:

Proto:
```proto
// addbook
// http POST /cs/v1/modname/add
rpc Add(AddReq) returns (AddResp) {}
```

API:
```text
@handler AddHandler
post /cs/v1/modname/add (addReq) returns (addResp)
```

### 4. Compatibility and safety
- Do not remove existing stable endpoints unless explicitly requested
- Regeneration must not silently break existing API contracts
- Clearly call out breaking changes when proto updates require them

## Execution workflow

1. Read and parse proto services/messages
2. Extract endpoint mapping source (annotation or convention)
3. Generate go-zero `type` blocks from messages
4. Generate endpoint declarations from RPCs
5. Validate naming consistency and contract compatibility
6. Return final `.api` plus a short mapping summary

## Clarification triggers

Stop and ask before generation when:
- Proto contains RPCs without route/method mapping source
- Multiple naming conventions are possible
- Backward compatibility constraints are not stated
- User expects partial generation but target scope is unclear

## Quality checklist

- Every RPC has method + route + request + response
- Every referenced type exists in generated `.api`
- No duplicate route signatures are introduced
- Existing contract impact is explicitly noted
- Generated content is deterministic and style-consistent
- Only service-related messages are converted
- HTTP comment method/path mapping is exact and case-normalized for go-zero keywords
- Request tags include both `json` and `form`

## Minimal output template

```text
syntax = "v1"

type (
  // mapped message types
)

@server(
  prefix: /v1
)
service xxx-api {
  @handler Xxx
  post /xxx (XxxReq) returns (XxxResp)
}
```
