# CURSOR 使用说明

在 Cursor 会话中需要把 `.proto` 转换为 go-zero `.api` 时，使用 `proto-to-gozero-api`。

## 推荐指令

```text
Use proto-to-gozero-api for this task.
Generate API definitions from the provided proto file.
If method/route mapping is not explicit, ask for the mapping policy first.
Keep compatibility unless the requirement explicitly asks for breaking changes.
```

## 检查清单

- service/rpc 解析完整
- message 到 type 映射完整
- 每个接口均有 method + route
- 兼容性影响有明确说明
