# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A custom status line utility for Claude Code written in Go. Displays model info, Git status, API usage limits, token consumption, and cost metrics in a 3-line formatted output with ANSI colors. Ships with 70+ themes selectable via interactive menu.

## Build Commands

```bash
# Build the binary (macOS / Linux)
go build -o statusline .

# Build on Windows
go build -o statusline.exe .

# Run (expects JSON input on stdin from Claude Code)
echo '{"model":{"display_name":"Claude Sonnet 4"},...}' | ./statusline

# Interactive theme menu
./statusline --menu

# Show version info
./statusline --version
```

## Architecture

**Single-file core**: Main logic in `statusline.go` (~1300 lines); themes live in `themes/`.

**Concurrency model**: 6 goroutines run in parallel to minimize latency (see `collectData`):
1. Git info (branch, staged/dirty counts)
2. Total hours calculation
3. Session usage parsing (token counts & cost from transcript)
4. Weekly stats loading
5. Daily stats loading
6. API usage fetching (file-cached for 30s)

The context window info no longer needs its own goroutine — Claude Code now passes it directly in the input JSON (`context_window` field).

Results are collected via a buffered channel and synchronized with `sync.WaitGroup`.

**Data flow**:
```
JSON Input (stdin) → Parse → Launch 6 parallel goroutines → Collect via channels
→ Update session/stats files → Format 3 output lines → Print to stdout
```

**Key data structures**:
- `Input`: Configuration from Claude Code (model, session ID, workspace, transcript path, context_window)
- `Config`: User config (theme, `usage_api` mode: `oauth_usage` default, or `haiku_probe`)
- `Session`: Tracks session intervals and total time
- `SessionUsageResult`: Token counts and cost for current session
- `APIUsage`: 5-hour and 7-day usage limits from Anthropic API

**External integrations**:
- OAuth token: reads `~/.claude/.credentials.json` first; falls back to macOS Keychain (`security find-generic-password -s "Claude Code-credentials"`)
- Anthropic API: `GET https://api.anthropic.com/api/oauth/usage` (default `oauth_usage` mode), or Haiku probe via `x-api-key` header (alternate `haiku_probe` mode)
- Git: Branch and status via `git branch --show-current` and `git status --porcelain`

**Config path** (`getConfigPath`): prefers XDG (`$XDG_CONFIG_HOME/claude-statusline/config.json` or `~/.config/claude-statusline/config.json`); falls back to binary-adjacent `config.json` for migration.

**Local storage** (in `~/.claude/session-tracker/`):
- `sessions/*.json`: Per-session time tracking
- `stats/daily-*.json` and `stats/weekly-*.json`: Accumulated statistics
- `api-usage-cache.json`: 30s-cached API usage response

## Model Pricing

Hardcoded in `modelPricing` map (per 1M tokens, current Claude 4.x family):

| Model  | Input | Output | Cache Read | Cache Write (5m) |
|--------|-------|--------|------------|------------------|
| Opus   | $5    | $25    | $0.50      | $6.25            |
| Sonnet | $3    | $15    | $0.30      | $3.75            |
| Haiku  | $1    | $5     | $0.10      | $1.25            |

Lookup is by model family (Opus/Sonnet/Haiku) via `getModelType` parsing `display_name`. Unknown models fall back to Sonnet pricing.

## Platform Support

Cross-platform: macOS, Linux, Windows. Terminal raw mode for `--menu` uses `golang.org/x/term`; Windows `\r\n` Enter handling is explicit in the menu loop.
