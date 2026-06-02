# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Accurate 1M/200K context window display using Claude Code's `context_window` JSON field
- Read 5h/7d rate limits directly from Claude Code's `rate_limits` JSON field (Pro/Max subscribers, recent versions), skipping the network round-trip / Haiku probe when present
- Lines-changed badge (`+added / -removed`) in the `twoline_pills` theme, from Claude Code's `cost.total_lines_added/removed`
- Fast mode awareness: detect `*-fast` model ids and apply ~2x rates in the fallback cost estimate

### Changed
- API usage now prefers the `rate_limits` stdin field and only falls back to the OAuth usage endpoint (or Haiku probe) when it is absent
- Session cost now prefers Claude Code's authoritative `cost.total_cost_usd` (covers Fast mode / future pricing) and only falls back to the transcript-based estimate when absent — removes most pricing-map maintenance
- Model pricing family is resolved from the stable `model.id` (falling back to `display_name`)
- Claude Code version is read from the `version` JSON field instead of shelling out to `claude --version` on every render (subprocess kept as fallback)
- Pricing comments updated to cover Opus 4.5–4.8 / Sonnet 4–4.6
- Documentation: `usage_api` default clarified as `oauth_usage`, with notes on `haiku_probe` fallback

### Dependencies
- Bump `golang.org/x/term` 0.40.0 → 0.41.0 → 0.43.0
- Bump `softprops/action-gh-release` 2 → 3
- Bump `codecov/codecov-action` 5 → 6

## [1.1.0] - 2026-03-11

### Added
- Linux support: terminal raw mode via `golang.org/x/term` (no more macOS-only)
- XDG config path (`$XDG_CONFIG_HOME/claude-statusline/config.json`) with binary-adjacent fallback for migration
- Configurable usage API mode (`usage_api`): `haiku_probe` (default at the time) or `oauth_usage`
- Haiku probe rate-limit detection as a fallback when the OAuth usage endpoint is persistently 429'd
- File-based cache for Haiku probe results
- Support for both epoch and ISO 8601 reset-time formats

### Fixed
- Windows `--menu` Enter key not working (handle `\r\n` pair, not only single `\r`/`\n`)
- Brighter pills themes; config path resolution for `--menu`/`--set-theme`
- Hardcoded model name now uses actual `display_name` from input
- Haiku probe uses `x-api-key` header (not Bearer)

### Changed
- CI: upgrade Go 1.21 → 1.24, golangci-lint v1 → v2 (config migration, lint rule cleanup)
- Release workflow Go version aligned with toolchain (1.24)

### Dependencies
- Bump `actions/upload-artifact` 4 → 6 → 7
- Bump `actions/download-artifact` 4 → 7 → 8
- Bump `actions/attest-build-provenance` 1 → 3 → 4
- Bump `actions/checkout` 4 → 6
- Bump `actions/setup-go` 5 → 6
- Bump `codecov/codecov-action` 4 → 5
- Bump `golangci/golangci-lint-action` 6 → 9
- Bump `goreleaser/goreleaser-action` 6 → 7
- Bump SLSA framework generator workflow to 2.1.0

## [1.0.3] - 2026-02-05

### Added
- 40 Japanese anime-themed statuslines (Shonen, Classic, Ghibli, etc.)
- Theme gallery pages; README simplified theme sections
- Theme screenshots and intro GIF
- EditorConfig, Dependabot, pkg.go.dev badge
- Makefile

### Changed
- Redesigned classic anime themes with unique visual styles
- Smoother intro.gif crossfade transitions
- `docs/` removed from version tracking (added to `.gitignore`)

### Fixed
- gofmt formatting in anime and ghibli theme files
- intro.gif properly animated
- Go module path corrected

## [1.0.2] - 2026-02-04

### Added
- GitHub Actions CI/CD pipeline
- Cross-platform release builds (Linux, macOS, Windows)
- Cosign artifact signing
- SBOM generation (SPDX format)
- SLSA Level 3 provenance
- `--version` flag to display version information
- Unit tests for core functions

### Changed
- Improved code formatting with `gofmt -s`

## [1.0.1] - 2026-02-04

### Added
- Initial CI/CD setup (superseded by v1.0.2)

## [1.0.0] - 2026-02-04

### Added
- Initial release
- Multi-theme support with 20+ themes
- Interactive theme selector (`--menu`)
- Real-time API usage monitoring
- Token consumption tracking
- Cost calculation with cache awareness
- Git status integration
- Session time tracking
- Daily/weekly/monthly statistics
