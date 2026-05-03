# Program-Gozero GitHub Standardization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 `program-gozero` 打造成专业、可上手、可协作的 GitHub 项目，补齐标准文档与仓库元信息。

**Architecture:** 采用“根仓库门户文档 + skills 子目录内容复用 + GitHub 协作模板”结构。根 README 提供中英双语入口、安装和调用方式；`.github` 提供 issue/PR 规范；仓库标准文件补齐许可、贡献、安全和变更记录。

**Tech Stack:** Markdown, GitHub Community Health Files, EditorConfig

---

### Task 1: 盘点当前仓库并确定标准化目标文件

**Files:**
- Create: `docs/superpowers/plans/2026-05-02-github-standardization.md`
- Modify: `README.md`（若不存在则创建）
- Create: `.github/*` 与根目录标准文件

- [ ] **Step 1: 确定目标文件列表**

目标文件：
- `README.md`
- `LICENSE` (MIT)
- `.gitignore`
- `.editorconfig`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `SECURITY.md`
- `CHANGELOG.md`
- `.github/ISSUE_TEMPLATE/bug_report.md`
- `.github/ISSUE_TEMPLATE/feature_request.md`
- `.github/pull_request_template.md`

- [ ] **Step 2: 校验是否与用户确认一致**

约束：
- README 为中英双语
- 采用标准版，不上 CI 重配置

### Task 2: 编写根 README（中英双语）

**Files:**
- Create/Modify: `README.md`

- [ ] **Step 1: 编写中文部分（优先）**

内容包含：
- 项目定位
- Skills 列表与用途
- 快速开始（安装/调用）
- 示例指令
- 仓库结构
- 贡献与许可证入口

- [ ] **Step 2: 编写英文对应部分**

保证中英文结构一致，降低维护难度。

### Task 3: 补齐 GitHub 标准文件

**Files:**
- Create: `LICENSE`, `.gitignore`, `.editorconfig`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, `CHANGELOG.md`
- Create: `.github/ISSUE_TEMPLATE/bug_report.md`, `.github/ISSUE_TEMPLATE/feature_request.md`, `.github/pull_request_template.md`

- [ ] **Step 1: 生成文件并填充可用模板**
- [ ] **Step 2: 对齐本项目场景（skills 仓库）**

### Task 4: 自检与交付

**Files:**
- Verify: 上述所有新增文件

- [ ] **Step 1: 检查链接、路径、命名与语义一致性**
- [ ] **Step 2: 运行基础状态检查（如 `git status`）**
- [ ] **Step 3: 形成最终变更说明并给出下一步建议**
