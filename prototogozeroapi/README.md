# prototogozeroapi

将 protobuf `.proto` 文件转换为 go-zero `.api` 文件的命令行工具。

## 功能说明

- 解析输入 `.proto` 文件中的所有结构定义
- 解析 `service` 和 `rpc` 定义
- 支持通过 `import` 解析外部 `.proto` 引用
- 优先使用注释中的 HTTP 映射：`// http POST /path`
- 若缺少 HTTP 映射，自动使用默认规则生成路由
- 在源 `.proto` 同目录下生成并覆盖同名 `.api` 文件

## 安装

在 `prototogozeroapi` 目录下执行：

```bash
go install ./cmd/prototogozeroapi
```

或者在仓库根目录执行：

```bash
go install ./prototogozeroapi/cmd/prototogozeroapi
```

## 使用方法

基本用法：

```bash
prototogozeroapi your.proto
```

带 import 搜索路径：

```bash
prototogozeroapi -I /path/to/includes -I /path/to/other your.proto
```

## 输出规则

- 输入：`demo.proto`
- 输出：`demo.api`
- 输出位置：与 `demo.proto` 相同目录

例如：

```bash
prototogozeroapi /data/proto/user.proto
```

生成结果：

```text
/data/proto/user.api
```

## 路由生成规则

优先读取 RPC 上方注释：

```proto
// http POST /cs/v1/user/add
rpc Add(AddReq) returns (AddResp) {}
```

生成：

```text
@handler AddHandler
post /cs/v1/user/add (AddReq) returns (AddResp)
```

如果没有 HTTP 注释，则使用默认规则：

```text
post /<packagePath>/<service>/<rpc>
```

例如：

- package: `cs.v1`
- service: `Book`
- rpc: `Add`

生成默认路由：

```text
post /cs/v1/book/add
```

## 注意事项

- 命令会直接覆盖目标 `.api` 文件
- `-I` 可以重复传入多个目录
- 如果 import 文件找不到，命令会直接返回错误
- 当前请求/响应字段会输出 `json` 和 `form` 标签，并附带 `,optional`

