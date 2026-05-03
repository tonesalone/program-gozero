# CLAUDE 使用说明

当需要把 `.proto` 合约转换为 go-zero `.api` 规范时，启用 `proto-to-gozero-api`。

## 推荐提示词

```text
Use proto-to-gozero-api.
Read <source>.proto and generate <target>.api.
Use proto HTTP annotations first.
If mapping is ambiguous, stop and ask before generating.
Do not introduce breaking route or type changes unless explicitly requested.
```

## 预期行为

- 从 proto 到 `.api` 的映射可复现、稳定
- method/route 规则不明确时先澄清再生成
- 对潜在破坏性变更给出明确兼容性说明
