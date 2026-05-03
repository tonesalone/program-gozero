# program-gozero

go-zero related agent skills collection.  
面向 `go-zero` 开发场景的 AI Skill 集合仓库。

## 中文说明

### 项目简介

`program-gozero` 用于集中维护可复用的 AI Skills，帮助你在 go-zero 相关开发中提升一致性与效率。  
当前已包含 2 个可直接使用的 skills：

- `program-doanot-skills`: 通用编码纪律与防回归守则
- `proto-to-gozero-api`: 从 `.proto` 生成 go-zero `.api` 的转换技能

### 快速开始

1. 克隆仓库：

```bash
git clone <your-repo-url> program-gozero
cd program-gozero
```

2. 将 `skills/` 下的目标 skill 复制到你的 AI 客户端技能目录（按你的客户端规范放置）。

3. 在对话中显式调用 skill，例如：

```text
Use program-doanot-skills for this task.
```

或：

```text
Use proto-to-gozero-api.
Generate user.api from user.proto.
Prefer proto HTTP annotations; if missing, ask for fallback route policy first.
```

### Skills 列表

| Skill | 作用 | 文档 |
|---|---|---|
| `program-doanot-skills` | 约束编码行为，强调最小改动、防回归、复杂度与终止条件 | `skills/program-doanot-skills/README.md` |
| `proto-to-gozero-api` | 将 proto 契约转换为 go-zero `.api` 规范 | `skills/proto-to-gozero-api/README.md` |

### 仓库结构

```text
program-gozero/
├── skills/
│   ├── program-doanot-skills/
│   └── proto-to-gozero-api/
├── .github/
└── README.md
```

### 常见问题

- Q: 为什么调用 skill 没有生效？  
  A: 请确认你使用的客户端已正确加载 skills 目录，并在提示词中显式调用 skill 名称。

- Q: 应该看中文还是英文文档？  
  A: 优先看与你协作环境一致的语言版本；若冲突，以 `SKILL.md` 的实际约束为准。

### 贡献

欢迎提交 Issue / PR。请先阅读 `CONTRIBUTING.md`。  
新增 skill 建议同时提供英文与中文文档，保持结构一致。

### 许可证

本项目采用 MIT License，详见 `LICENSE`。

---

## English

### Overview

`program-gozero` is a curated skill repository for AI-assisted go-zero development workflows.  
It currently includes two ready-to-use skills:

- `program-doanot-skills`: coding discipline and anti-regression guardrails
- `proto-to-gozero-api`: convert `.proto` contracts into go-zero `.api` specs

### Quick Start

1. Clone this repository:

```bash
git clone <your-repo-url> program-gozero
cd program-gozero
```

2. Copy target directories under `skills/` into your AI client's skill directory.

3. Invoke a skill explicitly in your prompt:

```text
Use program-doanot-skills for this task.
```

or:

```text
Use proto-to-gozero-api.
Generate user.api from user.proto.
Prefer proto HTTP annotations; if missing, ask for fallback route policy first.
```

### Repository Structure

```text
program-gozero/
├── skills/
│   ├── program-doanot-skills/
│   └── proto-to-gozero-api/
├── .github/
└── README.md
```

### Contributing

Please read `CONTRIBUTING.md` before opening issues or pull requests.

### License

Released under the MIT License. See `LICENSE`.
