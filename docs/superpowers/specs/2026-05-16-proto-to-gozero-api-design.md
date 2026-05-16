# Proto To go-zero API 设计说明

**目标**

实现一个可通过 `go install` 安装的命令行工具：输入一个 `.proto` 文件，解析其所有结构定义与服务 RPC，并在同目录生成同名 `.api` 文件（全量覆盖输出）。

**非目标**

- 不做增量对比、不读取历史基线；输出始终以最新 `.proto` 为准覆盖旧 `.api`。
- 不保证生成的 `.api` 与现有 `.api` 的稳定兼容（按用户要求忽略历史删除项等）。

---

## 输入/输出

- 输入：一个 `.proto` 文件路径（必填）
- 输出：默认输出到输入 `.proto` 的同目录
  - `foo.proto` → `foo.api`
- 依赖外部引用：支持解析 `import "xxx.proto";`，并解析被引用类型（通过 `-I` 指定 import 搜索路径，默认包含输入文件所在目录）

---

## 解析与范围

### 解析范围

- 对输入 `.proto` 文件：
  - 解析文件中所有 `message` 定义（无论是否被 RPC 引用）
  - 解析文件中所有 `service`/`rpc` 定义
- 对外部文件（通过 import 引入）：
  - 解析被引用到的结构定义（message）的闭包集合：只要被输入文件的 `message`/`rpc`/其传递字段引用到，就纳入生成

### 类型收集策略（确定性）

1. 收集输入文件内所有 `message`
2. 从以下入口做“传递引用闭包”收集额外 `message`：
   - 所有 `rpc` 的 request/response 类型
   - 输入文件内所有 `message` 的字段类型
3. 若闭包中出现外部文件 `message`，继续解析并递归收集其字段引用的 `message`

> 说明：`enum` 会被解析用于字段类型映射，但 go-zero `.api` 不支持 `type alias`（因此不输出 enum 类型本体），enum 字段统一映射为 `int32`。

---

## 映射规则

### 1) message → type（go-zero struct）

- type 名：保持 proto message 名不变
  - 例：`AddReq` → `AddReq`
- 字段顺序：保持 proto 中声明顺序
- 字段名：保持 proto 字段名不变
- request/response tag：
  - request 类型字段：同时包含 `json` 与 `form`，并且两者都必须带 `,optional`
  - response 类型字段：至少包含 `json`，并且必须带 `,optional`
  - 例：`book string \`json:"book,optional" form:"book,optional"\``

### 2) 字段类型映射（protobuf → go-zero api 类型）

- 标量：
  - `string` → `string`
  - `bool` → `bool`
  - `bytes` → `string`
  - `int32/sint32/sfixed32` → `int32`
  - `uint32/fixed32` → `uint32`
  - `int64/sint64/sfixed64` → `int64`
  - `uint64/fixed64` → `uint64`
  - `float` → `float32`
  - `double` → `float64`
- `enum`：`int32`
- `repeated T`：`[]T`
- `map<K,V>`：`map[K]V`（K/V 均使用上述映射后的类型）
- message 引用：映射到对应的 type 名

### 3) rpc → endpoint

- 一个 `rpc` 生成一个 endpoint
- handler 命名：`@handler <RpcName>Handler`
- 请求/响应类型：严格与 `rpc` 签名对齐

### 4) 路由与方法

优先使用 proto 注释中的显式 HTTP 标注，格式：

`// http METHOD /path`

- `METHOD` 取值如：`GET/POST/PUT/DELETE/HEAD/PATCH`
- 输出时 method 关键字小写：`get/post/...`
- path 原样使用（不改写）

如果某个 `rpc` 缺少 HTTP 标注，使用内置默认规则（全文件统一，不混用）：

- method：`post`
- path：`/` + `<packagePath>` + `/` + `<serviceName>` + `/` + `<rpcName>` 全部转小写
  - `<packagePath>`：proto `package a.b;` → `a/b`，如果无 package 则为空（省略该段）

---

## 输出 `.api` 模板

生成文件结构固定：

- `syntax = "v1"`
- `type ( ... )`：包含收集到的所有 message 类型
- `service ... { ... }`：对每个 proto service 生成一个 service block
  - service 名：`strings.ToLower(<ServiceName>) + "-api"`

---

## 冲突与边界处理

- 嵌套 message：`.api` 没有作用域，若出现同名冲突，采用确定性扁平化命名：
  - `Outer.Inner` → `Outer_Inner`
- 跨文件同名 message：若发生冲突，采用确定性前缀：
  - `<package>_<TypeName>`（package 中 `.` 替换为 `_`）
- 未支持/未知语法：遇到无法解析/无法映射的字段类型时，工具直接失败并返回错误信息（不输出半成品文件）

