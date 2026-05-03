# proto-to-gozero-api

用于将 protobuf `.proto` 定义转换为 go-zero `.api` 文件。

## 功能说明

- 将 protobuf `message` 映射为 go-zero `type`
- 将 protobuf `rpc` 映射为 go-zero 接口声明
- 依据注解或约定生成 method/route
- 在未明确要求变更时优先保持兼容

## 使用场景

- 首次根据 proto 合约生成 go-zero API 规范
- proto 更新后重新生成 `.api`
- 团队内统一 proto 到 API 的转换规则

## 文件说明

- `SKILL.md`：英文技能定义
- `SKILL.zh.md`：中文技能定义
- `CLAUDE.md`：Claude 使用说明
- `CURSOR.md`：Cursor 使用说明

## 快速用法

可以这样给代理下指令：

```text
Use proto-to-gozero-api.
Generate <target>.api from <source>.proto.
Prefer proto HTTP annotations; if missing, ask for fallback route policy before generation.
```
