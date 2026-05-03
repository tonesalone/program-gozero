---
name: "proto-to-gozero-api"
description: "从 protobuf .proto 生成 go-zero .api。用户要求将 proto 的 service/rpc/message 转为 go-zero 路由、类型与接口声明时调用。"
---

# 从 Proto 生成 go-zero API

将 protobuf `.proto` 定义稳定、可复现地转换为 go-zero `.api` 文件。

## 何时调用

在以下场景调用本技能：
- 用户要求根据 `.proto` 生成 go-zero `.api`
- 需要把 proto 的 service/rpc/message 映射成 go-zero API 结构
- 需要从 proto 推导接口方法、路由、请求和响应类型
- `.proto` 更新后需要同步再生成 `.api`

## 必要输入

- 源 `.proto` 文件路径
- 目标 `.api` 文件路径
- 路由策略：
  - 优先使用 proto 中显式 HTTP 注解（若存在）
  - 否则使用约定的兜底映射规则（必须明确）
- 命名策略：
  - 保持 proto 原命名，或
  - 转换为项目命名风格（必须明确）

## 输出要求

生成结果应包含：
- 由 protobuf message 映射出的 `type` 块（仅限 service 接口涉及）
- 由 RPC 映射出的接口声明
- 每个接口的 HTTP method + route path
- 与 proto 对齐的请求/响应类型引用

严格范围规则：
- 只转换 proto `service` 域中接口所涉及的数据结构
- 未被 service 的 RPC 请求/响应链路引用的 message，不转换

## 核心映射规则

### 1. Message 到 type
- protobuf `message` 映射为 go-zero `type`
- 除非项目规范要求，否则保持字段顺序
- 保留语义化字段名，避免无必要重命名
- 类型命名遵循样式：`AddReq` -> `addReq`，`CheckResp` -> `checkResp`
- 请求结构必须包含 `json` 和 `form` 两类 tag
- 响应结构至少包含 `json` tag

示例：

Proto：
```proto
message AddReq {
  string book = 1;
  int64 price = 2;
}
```

API：
```text
addReq {
  book string `json:"book" form:"book"`
  price int64 `json:"price" form:"price"`
}
```

### 2. RPC 到接口
- protobuf `rpc` 映射为一个 API 处理入口
- 请求/响应类型必须与 proto RPC 签名一致
- 若 method/route 缺失，先询问映射规则再生成
- Handler 命名遵循 RPC 名称：`Add` -> `@handler AddHandler`

### 3. 路由与方法策略
- 第一优先级：使用 proto 中显式 HTTP 注释
- 兜底策略：统一应用团队/项目约定
- 同一文件中禁止混用未确认的多套约定
- 解析注释样式：`// http POST /cs/v1/modname/add`
- 将 `POST` 映射为 `post`，`GET` 映射为 `get`
- 路径与注释保持一致
- 路径中的模块名（如 `modname`）必须与 proto 的模块/包语义一致

示例：

Proto：
```proto
// 添加书籍
// http POST /cs/v1/modname/add
rpc Add(AddReq) returns (AddResp) {}
```

API：
```text
@handler AddHandler
post /cs/v1/modname/add (addReq) returns (addResp)
```

### 4. 兼容性与安全
- 除非需求明确要求，不删除已有稳定接口
- 重新生成不得静默破坏现有 API 合约
- 对必须发生的破坏性变更要明确标注

## 执行流程

1. 解析 proto 的 service/message/rpc
2. 提取接口映射来源（注解或约定）
3. 生成 go-zero `type` 块
4. 生成 RPC 对应接口声明
5. 校验命名一致性与兼容性
6. 输出 `.api` 与简短映射说明

## 必须先澄清的情况

遇到以下情况先提问，不直接生成：
- RPC 没有 method/route 来源
- 命名规则存在多种可能
- 回归兼容要求未明确
- 只需部分生成但范围不清楚

## 质量检查清单

- 每个 RPC 都有 method + route + req + resp
- 每个引用类型都在 `.api` 中定义
- 不引入重复路由签名
- 兼容性影响有明确说明
- 生成结果可复现且风格一致
- 只转换 service 相关的数据结构
- HTTP 注释中的 method/path 与 API 映射准确一致
- 请求结构同时具备 `json` 与 `form` tag
