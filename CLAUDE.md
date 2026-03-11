# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A custom status line utility for Claude Code written in Go. Displays model info, Git status, API usage limits, token consumption, and cost metrics in a 3-line formatted output with ANSI colors.

## Build Commands

```bash
# Build the binary
go build -o statusline statusline.go

# Build on Windows
go build -o statusline.exe .

# Run (expects JSON input on stdin from Claude Code)
echo '{"model":{"display_name":"Claude Sonnet 4"},...}' | ./statusline
```

## Architecture

**Single-file design**: All code is in `statusline.go` (~1100 lines).

**Concurrency model**: 7 goroutines run in parallel to minimize latency:
1. Git info (branch, staged/dirty counts)
2. Total hours calculation
3. Context window analysis
4. Session usage parsing
5. Weekly stats loading
6. Daily stats loading
7. API usage fetching (with 30s cache)

Results are collected via channels and synchronized with `sync.WaitGroup`.

**Data flow**:
```
JSON Input (stdin) → Parse → Launch 7 parallel goroutines → Collect via channels
→ Update session/stats files → Format 3 output lines → Print to stdout
```

**Key data structures**:
- `Input`: Configuration from Claude Code (model, session ID, workspace, transcript path)
- `Session`: Tracks session intervals and total time
- `SessionUsageResult`: Token counts and cost for current session
- `APIUsage`: 5-hour and 7-day usage limits from Anthropic API

**External integrations**:
- macOS Keychain: Reads OAuth token via `security find-generic-password`
- Anthropic API: `GET https://api.anthropic.com/api/oauth/usage`
- Git: Branch and status via `git branch --show-current` and `git status --porcelain`

**Local storage** (in `~/.claude/session-tracker/`):
- `sessions/*.json`: Per-session time tracking
- `stats/daily-*.json` and `stats/weekly-*.json`: Accumulated statistics

## Model Pricing

Hardcoded in `modelPricing` map (per 1M tokens):
- Opus: $15 input, $75 output
- Sonnet: $3 input, $15 output
- Haiku: $0.25 input, $1.25 output

## Platform Requirement

macOS only - depends on Keychain for OAuth token access.
