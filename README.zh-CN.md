# Claude Statusline

[![CI](https://github.com/kevinlincg/claude-statusline/actions/workflows/ci.yml/badge.svg)](https://github.com/kevinlincg/claude-statusline/actions/workflows/ci.yml)
[![Release](https://github.com/kevinlincg/claude-statusline/actions/workflows/release.yml/badge.svg)](https://github.com/kevinlincg/claude-statusline/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kevinlincg/claude-statusline)](https://goreportcard.com/report/github.com/kevinlincg/claude-statusline)
[![Go Reference](https://pkg.go.dev/badge/github.com/kevinlincg/claude-statusline.svg)](https://pkg.go.dev/github.com/kevinlincg/claude-statusline)
[![GitHub release](https://img.shields.io/github/v/release/kevinlincg/claude-statusline)](https://github.com/kevinlincg/claude-statusline/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)

[English](README.md) | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md) | [日本語](README.ja.md)

为 Claude Code 打造的自定义状态栏，使用 Go 语言编写。显示模型信息、Git 状态、API 使用量、Token 消耗、成本指标等。

<p align="center">
  <img src="assets/intro.gif" alt="Claude Statusline 主题展示" width="750">
</p>

## 安装

### 方式一：下载可执行文件（推荐）

从 [GitHub Releases](https://github.com/kevinlincg/claude-statusline/releases/latest) 下载适合您平台的版本：

```bash
# macOS (Apple Silicon)
curl -L https://github.com/kevinlincg/claude-statusline/releases/latest/download/claude-statusline_darwin_arm64.tar.gz | tar xz
mv statusline ~/.claude/statusline-go/

# macOS (Intel)
curl -L https://github.com/kevinlincg/claude-statusline/releases/latest/download/claude-statusline_darwin_amd64.tar.gz | tar xz
mv statusline ~/.claude/statusline-go/

# Linux (x64)
curl -L https://github.com/kevinlincg/claude-statusline/releases/latest/download/claude-statusline_linux_amd64.tar.gz | tar xz
mv statusline ~/.claude/statusline-go/

# Linux (ARM64)
curl -L https://github.com/kevinlincg/claude-statusline/releases/latest/download/claude-statusline_linux_arm64.tar.gz | tar xz
mv statusline ~/.claude/statusline-go/
```

### 方式二：从源码编译

```bash
# 克隆项目
git clone https://github.com/kevinlincg/claude-statusline.git ~/.claude/statusline-go

# 编译
cd ~/.claude/statusline-go
go build -o statusline .
```

### 配置 Claude Code

编辑 `~/.claude/settings.json`：

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/statusline-go/statusline"
  }
}
```

## 主题

### 交互式主题选择器

使用交互式菜单预览并选择主题：

```bash
./statusline --menu
```

使用方向键（或 h/l）浏览主题，Enter 确认，q 取消。

### 命令行选项

```bash
./statusline --list-themes      # 列出所有主题
./statusline --preview <name>   # 预览特定主题
./statusline --set-theme <name> # 直接设置主题
./statusline --menu             # 交互式主题选择器
./statusline --version          # 显示版本信息
```

### 手动配置

编辑 `~/.claude/statusline-go/config.json`：

```json
{
  "theme": "classic_framed"
}
```

### 可用主题

目前提供 **65 种主题**：

| 分类 | 主题 |
|------|------|
| 经典 & 简约 | `classic`, `classic_framed`, `minimal`, `compact`, `boxed`, `zen` |
| 科幻 & 赛博朋克 | `hud`, `cyberpunk`, `synthwave`, `matrix`, `glitch` |
| 自然 & 美学 | `ocean`, `steampunk` |
| 系统监视器 | `htop`, `btop`, `gtop`, `stui` |
| 复古 & 游戏 | `pixel`, `retro_crt`, `bbs`, `lord`, `tradewars`, `nethack`, `dungeon`, `mud_rpg` |

**[查看所有主题截图 →](THEMES.zh-CN.md)**

## 显示信息

### 第一行：基本信息
- **模型**：当前使用的 Claude 模型（Opus/Sonnet/Haiku）
- **项目**：当前工作目录名称
- **Git 分支**：分支名称与状态（+已暂存/~未暂存）
- **Context**：Context Window 使用量进度条
- **每日工时**：今日累计工作时间

### 第二行：API 限制
- **Session**：5 小时内 API 使用率与重置时间
- **Week**：7 天内 API 使用率与重置时间

进度条颜色：绿色 (<50%) → 黄色 (50-75%) → 橙色 (75-90%) → 红色 (>90%)

### 第三行：Session 统计
- **Token**：本次 Session 累计使用的 Token 数量
- **成本**：本次 Session 的预估成本 (USD)
- **时长**：Session 持续时间
- **消息数**：对话消息数量
- **烧钱速度**：每小时花费
- **今日/周成本**：累计成本
- **Cache 命中率**：Cache read 比例（绿色 ≥70% / 黄色 40-70% / 橙色 <40%）

## 定价

每百万 Token（2026 年 1 月）：

| 模型 | 输入 | 输出 | Cache 读取 | Cache 写入 |
|------|------|------|------------|------------|
| Opus 4.5 | $5 | $25 | $0.50 | $6.25 |
| Sonnet 4/4.5 | $3 | $15 | $0.30 | $3.75 |
| Haiku 4.5 | $1 | $5 | $0.10 | $1.25 |

## 数据存储

统计数据保存于 `~/.claude/session-tracker/`：
- `sessions/` - 单个 Session 数据
- `stats/` - 每日与每周 Token 统计

## 贡献

欢迎贡献！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

## 安全性

发布的文件均经过签名并包含 SLSA 来源证明。请参阅 [SECURITY.md](SECURITY.md) 了解验证方式。

## 许可证

[MIT License](LICENSE)
