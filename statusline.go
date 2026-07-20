package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kevinlincg/claude-statusline/themes"
	"golang.org/x/term"
)

// Version information (set via ldflags during build)
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// Model pricing (per 1M tokens)
var modelPricing = map[string]struct {
	Input      float64
	Output     float64
	CacheRead  float64
	CacheWrite float64
}{
	"Opus": {
		Input:      5.0,  // Opus 4.5–4.8: $5 per 1M input tokens
		Output:     25.0, // Opus 4.5–4.8: $25 per 1M output tokens
		CacheRead:  0.5,  // Opus 4.5–4.8: $0.50 per 1M cache read tokens
		CacheWrite: 6.25, // Opus 4.5–4.8: $6.25 per 1M cache write tokens (5m)
	},
	"Sonnet": {
		Input:      3.0,  // Sonnet 4–4.6: $3 per 1M input tokens
		Output:     15.0, // Sonnet 4–4.6: $15 per 1M output tokens
		CacheRead:  0.3,  // Sonnet 4–4.6: $0.30 per 1M cache read tokens
		CacheWrite: 3.75, // Sonnet 4–4.6: $3.75 per 1M cache write tokens (5m)
	},
	"Haiku": {
		Input:      1.0,  // Haiku 4.5: $1 per 1M input tokens
		Output:     5.0,  // Haiku 4.5: $5 per 1M output tokens
		CacheRead:  0.1,  // Haiku 4.5: $0.10 per 1M cache read tokens
		CacheWrite: 1.25, // Haiku 4.5: $1.25 per 1M cache write tokens (5m)
	},
}

// Input data structure
type Input struct {
	Model struct {
		DisplayName string `json:"display_name"`
		ID          string `json:"id"` // e.g. "claude-opus-4-8" / "claude-opus-4-8-fast"
	} `json:"model"`
	SessionID string `json:"session_id"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	// Version is the Claude Code version, supplied directly so we can avoid
	// shelling out to `claude --version` on every render.
	Version string `json:"version"`
	// Cost is computed client-side by Claude Code. total_cost_usd is the
	// authoritative session cost (covers Fast mode, batch, etc.); we prefer it
	// over our own transcript-based estimate when present.
	Cost struct {
		TotalCostUSD      float64 `json:"total_cost_usd"`
		TotalLinesAdded   int     `json:"total_lines_added"`
		TotalLinesRemoved int     `json:"total_lines_removed"`
	} `json:"cost"`
	TranscriptPath string `json:"transcript_path,omitempty"`
	ContextWindow  struct {
		ContextWindowSize int `json:"context_window_size"`
		TotalInputTokens  int `json:"total_input_tokens"`
		TotalOutputTokens int `json:"total_output_tokens"`
		UsedPercentage    int `json:"used_percentage"`
	} `json:"context_window"`
	// RateLimits is supplied directly by Claude Code (Pro/Max subscribers,
	// recent versions) so we can skip the network round-trip when present.
	// resets_at is Unix epoch seconds; used_percentage is 0-100.
	RateLimits struct {
		FiveHour struct {
			UsedPercentage float64 `json:"used_percentage"`
			ResetsAt       int64   `json:"resets_at"`
		} `json:"five_hour"`
		SevenDay struct {
			UsedPercentage float64 `json:"used_percentage"`
			ResetsAt       int64   `json:"resets_at"`
		} `json:"seven_day"`
	} `json:"rate_limits"`
}

// Config structure
type Config struct {
	Theme    string `json:"theme"`
	UsageAPI string `json:"usage_api,omitempty"` // "oauth_usage" (default) or "haiku_probe"
}

// Session data structure
type Session struct {
	ID            string     `json:"id"`
	Date          string     `json:"date"`
	Start         int64      `json:"start"`
	LastHeartbeat int64      `json:"last_heartbeat"`
	TotalSeconds  int64      `json:"total_seconds"`
	Intervals     []Interval `json:"intervals"`
}

type Interval struct {
	Start int64  `json:"start"`
	End   *int64 `json:"end"`
}

// UsageStats structure
type UsageStats struct {
	TotalCost    float64            `json:"total_cost"`
	SessionCosts map[string]float64 `json:"session_costs,omitempty"`
	Date         string             `json:"date"`
	Week         string             `json:"week"`
	LastUpdated  int64              `json:"last_updated"`
}

// APIUsage structure
type APIUsage struct {
	FiveHour struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"seven_day"`
}

// Result channel data
type Result struct {
	Type string
	Data interface{}
}

// GitInfo contains Git status information
type GitInfo struct {
	Branch      string
	DirtyCount  int
	StagedCount int
	AheadCount  int    // commits ahead of upstream (to push)
	BehindCount int    // commits behind upstream (to pull)
	StashCount  int    // number of stash entries
	ShortSHA    string // short commit SHA of HEAD
}

// SessionUsageResult contains session usage information
type SessionUsageResult struct {
	InputTokens      int64
	OutputTokens     int64
	CacheReadTokens  int64
	CacheWriteTokens int64
	Cost             float64
	MessageCount     int
	Duration         time.Duration
}

// APIUsageCache wraps APIUsage with a timestamp for file-based caching.
type APIUsageCache struct {
	Usage    APIUsage  `json:"usage"`
	CachedAt time.Time `json:"cached_at"`
}

func main() {
	// Command line arguments
	listThemes := flag.Bool("list-themes", false, "List all available themes")
	previewTheme := flag.String("preview", "", "Preview a specific theme")
	setTheme := flag.String("set-theme", "", "Set theme")
	menuMode := flag.Bool("menu", false, "Interactive theme menu")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Process command line arguments
	if *showVersion {
		fmt.Printf("statusline %s (commit: %s, built: %s)\n", Version, Commit, Date)
		return
	}

	if *listThemes {
		printThemeList()
		return
	}

	if *previewTheme != "" {
		previewThemeDemo(*previewTheme)
		return
	}

	if *setTheme != "" {
		saveThemeConfig(*setTheme)
		return
	}

	if *menuMode {
		runInteractiveMenu()
		return
	}

	// Normal mode: read stdin and output statusline
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode input: %v\n", err)
		os.Exit(1)
	}

	// Get model type (prefers the stable model.id over the display name)
	modelType := getModelTypeFromInput(input)

	// Collect data in parallel
	data := collectData(input, modelType)

	// Update session and stats
	updateSession(input.SessionID)
	updateDailyStats(input.SessionID, data, modelType)

	// Load theme config
	themeName := loadThemeConfig()
	theme, ok := themes.GetTheme(themeName)
	if !ok {
		// Default theme
		theme, _ = themes.GetTheme("classic_framed")
	}

	// Render output
	fmt.Print(theme.Render(data))
}

// printThemeList lists all available themes
func printThemeList() {
	fmt.Println("\nAvailable themes:")
	fmt.Println(strings.Repeat("─", 60))

	themeList := themes.ListThemes()
	sort.Slice(themeList, func(i, j int) bool {
		return themeList[i].Name() < themeList[j].Name()
	})

	for _, t := range themeList {
		fmt.Printf("  %-16s  %s\n", t.Name(), t.Description())
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Println("\nUsage:")
	fmt.Println("  ./statusline --set-theme <theme-name>  Set theme")
	fmt.Println("  ./statusline --preview <theme-name>    Preview theme")
	fmt.Println("  ./statusline --menu                    Interactive menu")
	fmt.Println()
}

// runInteractiveMenu runs interactive theme menu
func runInteractiveMenu() {
	themeList := themes.ListThemes()
	sort.Slice(themeList, func(i, j int) bool {
		return themeList[i].Name() < themeList[j].Name()
	})

	if len(themeList) == 0 {
		fmt.Println("No themes available")
		return
	}

	// Find current theme
	currentTheme := loadThemeConfig()
	selectedIndex := 0
	for i, t := range themeList {
		if t.Name() == currentTheme {
			selectedIndex = i
			break
		}
	}

	// Set terminal to raw mode
	oldState, err := makeRaw(os.Stdin.Fd())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set terminal: %v\n", err)
		return
	}
	defer restore(os.Stdin.Fd(), oldState)

	// Test data
	testData := themes.StatusData{
		ModelName:       "Opus 4.6",
		ModelType:       "Opus",
		Version:         "v1.0.75",
		UpdateAvailable: true,
		ProjectPath:     "~/cookys/project",
		GitBranch:       "main",
		GitStaged:       3,
		GitDirty:        5,
		GitAhead:        2,
		GitBehind:       1,
		GitStash:        1,
		GitSHA:          "a1b2c3d",
		TokensPerSec:    1250.0,
		TokenCount:      45200,
		MessageCount:    12,
		SessionTime:     "1h30m",
		CacheHitRate:    78,
		SessionCost:     0.12,
		DayCost:         3.45,
		MonthCost:       67.89,
		WeekCost:        23.45,
		BurnRate:        5.2,
		ContextUsed:     90000,
		ContextPercent:  45,
		API5hrPercent:   23,
		API5hrTimeLeft:  "3h17m",
		API7dayPercent:  67,
		API7dayTimeLeft: "2d5h",
		LinesAdded:      156,
		LinesRemoved:    23,
	}

	// Print function (raw mode requires \r\n)
	println := func(s string) {
		fmt.Print(s + "\r\n")
	}

	renderMenu := func() {
		// Clear screen
		fmt.Print("\033[2J\033[H")

		// Previous theme name
		prevName := ""
		if selectedIndex > 0 {
			prevName = themeList[selectedIndex-1].Name()
		}

		// Next theme name
		nextName := ""
		if selectedIndex < len(themeList)-1 {
			nextName = themeList[selectedIndex+1].Name()
		}

		// Title bar: show previous, current, next
		println(fmt.Sprintf("\033[1mTheme Selector\033[0m   \033[2m%12s <\033[0m \033[1;7m %s \033[0m \033[2m> %-12s\033[0m",
			prevName, themeList[selectedIndex].Name(), nextName))
		println(fmt.Sprintf("   %s", themeList[selectedIndex].Description()))
		println(strings.Repeat("─", 100))

		// Preview (replace \n with \r\n)
		preview := themeList[selectedIndex].Render(testData)
		preview = strings.ReplaceAll(preview, "\n", "\r\n")
		fmt.Print(preview)

		println(strings.Repeat("─", 100))
		println("\033[2m< > Select theme  |  Enter Confirm  |  q Cancel\033[0m")
	}

	renderMenu()

	// Read keypress
	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		// Handle Enter: single \r or \n (n==1), or \r\n pair (n==2, Windows)
		if (n == 1 && (buf[0] == 13 || buf[0] == 10)) ||
			(n == 2 && buf[0] == 13 && buf[1] == 10) {
			fmt.Print("\033[2J\033[H")
			saveThemeConfig(themeList[selectedIndex].Name())
			fmt.Printf("Theme set to: %s\r\n", themeList[selectedIndex].Name())
			return
		}

		if n == 1 {
			switch buf[0] {
			case 'q', 'Q', 27: // q or Escape
				fmt.Print("\033[2J\033[H")
				fmt.Print("Canceled\r\n")
				return
			case 'h', 'H': // vim-style left
				if selectedIndex > 0 {
					selectedIndex--
					renderMenu()
				}
			case 'l', 'L': // vim-style right
				if selectedIndex < len(themeList)-1 {
					selectedIndex++
					renderMenu()
				}
			}
		} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
			// Arrow keys
			switch buf[2] {
			case 68: // Left
				if selectedIndex > 0 {
					selectedIndex--
					renderMenu()
				}
			case 67: // Right
				if selectedIndex < len(themeList)-1 {
					selectedIndex++
					renderMenu()
				}
			}
		}
	}
}

// Terminal raw mode functions
func makeRaw(fd uintptr) (*term.State, error) {
	return term.MakeRaw(int(fd))
}

func restore(fd uintptr, oldState *term.State) {
	if oldState != nil {
		term.Restore(int(fd), oldState)
	}
}

// previewThemeDemo previews a theme
func previewThemeDemo(themeName string) {
	theme, ok := themes.GetTheme(themeName)
	if !ok {
		fmt.Printf("Error: theme '%s' not found\n", themeName)
		fmt.Println("Use --list-themes to see all available themes")
		return
	}

	// Create test data
	data := themes.StatusData{
		ModelName:       "Opus 4.6",
		ModelType:       "Opus",
		Version:         "v1.0.75",
		UpdateAvailable: true,
		ProjectPath:     "~/cookys/project",
		GitBranch:       "main",
		GitStaged:       3,
		GitDirty:        5,
		GitAhead:        2,
		GitBehind:       1,
		GitStash:        1,
		GitSHA:          "a1b2c3d",
		TokensPerSec:    1250.0,
		TokenCount:      45200,
		MessageCount:    12,
		SessionTime:     "1h30m",
		CacheHitRate:    78,
		SessionCost:     0.12,
		DayCost:         3.45,
		MonthCost:       67.89,
		WeekCost:        23.45,
		BurnRate:        5.2,
		ContextUsed:     90000,
		ContextPercent:  45,
		API5hrPercent:   23,
		API5hrTimeLeft:  "3h17m",
		API7dayPercent:  67,
		API7dayTimeLeft: "2d5h",
		LinesAdded:      156,
		LinesRemoved:    23,
	}

	fmt.Printf("\nPreview theme: %s\n", themeName)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Print(theme.Render(data))
	fmt.Println()
}

// saveThemeConfig saves theme configuration
func saveThemeConfig(themeName string) {
	// Check if theme exists
	if _, ok := themes.GetTheme(themeName); !ok {
		fmt.Printf("Error: theme '%s' not found\n", themeName)
		fmt.Println("Use --list-themes to see all available themes")
		return
	}

	configFile := getConfigPath()
	os.MkdirAll(filepath.Dir(configFile), 0755)

	config := loadConfig()
	config.Theme = themeName
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0644)

	fmt.Printf("Theme set to: %s\n", themeName)
}

// loadConfig loads the full configuration from file.
func loadConfig() Config {
	configFile := getConfigPath()
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}
	}
	var config Config
	json.Unmarshal(data, &config)
	return config
}

// loadThemeConfig loads theme configuration
func loadThemeConfig() string {
	config := loadConfig()
	if config.Theme == "" {
		return "classic_framed"
	}
	return config.Theme
}

// getConfigPath returns the config file path
// Priority: XDG config dir > binary-adjacent (migration fallback)
func getConfigPath() string {
	// 1. XDG / ~/.config location
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			configDir = filepath.Join(homeDir, ".config")
		}
	}
	if configDir != "" {
		xdgPath := filepath.Join(configDir, "claude-statusline", "config.json")
		if _, err := os.Stat(xdgPath); err == nil {
			return xdgPath
		}
		// Check binary-adjacent for migration
		exe, err := os.Executable()
		if err == nil {
			exe, _ = filepath.EvalSymlinks(exe)
			adjPath := filepath.Join(filepath.Dir(exe), "config.json")
			if _, err := os.Stat(adjPath); err == nil {
				return adjPath
			}
		}
		// Default to XDG path (will be created on first save)
		return xdgPath
	}
	// Ultimate fallback
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "claude-statusline", "config.json")
}

// collectData collects all data
func collectData(input Input, modelType string) themes.StatusData {
	results := make(chan Result, 10)
	var wg sync.WaitGroup

	wg.Add(6)

	go func() {
		defer wg.Done()
		gitInfo := getGitInfo()
		results <- Result{"git", gitInfo}
	}()

	go func() {
		defer wg.Done()
		totalHours := calculateTotalHours(input.SessionID)
		results <- Result{"hours", totalHours}
	}()

	fastMode := isFastMode(input)
	go func() {
		defer wg.Done()
		sessionUsage := calculateSessionUsage(input.TranscriptPath, input.SessionID, modelType, fastMode)
		results <- Result{"session_usage", sessionUsage}
	}()

	go func() {
		defer wg.Done()
		weeklyStats := getWeeklyStats()
		results <- Result{"weekly", weeklyStats}
	}()

	go func() {
		defer wg.Done()
		dailyStats := getDailyStats()
		results <- Result{"daily", dailyStats}
	}()

	go func() {
		defer wg.Done()
		// Prefer rate_limits supplied in Claude Code's JSON input; only fall
		// back to the network fetch (or Haiku probe) when it's absent.
		apiUsage := apiUsageFromInput(input)
		if apiUsage == nil {
			apiUsage = fetchAPIUsage()
		}
		results <- Result{"api_usage", apiUsage}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var (
		gitInfo      GitInfo
		totalHours   string
		sessionUsage SessionUsageResult
		dailyStats   UsageStats
		weeklyStats  UsageStats
		apiUsage     *APIUsage
	)

	for result := range results {
		switch result.Type {
		case "git":
			gitInfo = result.Data.(GitInfo)
		case "hours":
			totalHours = result.Data.(string)
		case "session_usage":
			sessionUsage = result.Data.(SessionUsageResult)
		case "weekly":
			weeklyStats = result.Data.(UsageStats)
		case "daily":
			dailyStats = result.Data.(UsageStats)
		case "api_usage":
			apiUsage = result.Data.(*APIUsage)
		}
	}

	// Context from Claude Code's JSON input (supports both 200K and 1M windows)
	contextUsed := input.ContextWindow.TotalInputTokens + input.ContextWindow.TotalOutputTokens
	contextPercent := input.ContextWindow.UsedPercentage

	// Get monthly stats
	monthlyStats := getMonthlyStats()

	// Calculate burn rate
	burnRate := calculateBurnRateValue(dailyStats)

	// Get version and update status (prefers input.Version, no subprocess)
	version, updateAvailable := resolveVersion(input.Version)

	// API data
	api5hrPercent := 0
	api5hrTimeLeft := "--"
	api7dayPercent := 0
	api7dayTimeLeft := "--"

	if apiUsage != nil {
		api5hrPercent = int(apiUsage.FiveHour.Utilization)
		api5hrTimeLeft = formatTimeLeftShort(apiUsage.FiveHour.ResetsAt)
		api7dayPercent = int(apiUsage.SevenDay.Utilization)
		api7dayTimeLeft = formatTimeLeftShort(apiUsage.SevenDay.ResetsAt)
	}

	// Calculate cache hit rate
	cacheHitRate := 0
	totalInput := sessionUsage.InputTokens + sessionUsage.CacheReadTokens
	if totalInput > 0 {
		cacheHitRate = int(float64(sessionUsage.CacheReadTokens) * 100.0 / float64(totalInput))
	}

	// Token throughput (tokens/second) over the session's active span. Guard the
	// zero-duration case (single message, or a transcript with no timestamps) to
	// avoid a divide-by-zero producing +Inf.
	totalTokens := sessionUsage.InputTokens + sessionUsage.OutputTokens + sessionUsage.CacheReadTokens + sessionUsage.CacheWriteTokens
	tokensPerSec := 0.0
	if secs := sessionUsage.Duration.Seconds(); secs > 0 {
		tokensPerSec = float64(totalTokens) / secs
	}

	// Session cost: prefer Claude Code's authoritative client-side value
	// (covers Fast mode, batch, and future pricing changes); fall back to our
	// transcript-based estimate when it's absent. The downstream daily/weekly/
	// monthly accumulators key off sessionID deltas, so the source can switch
	// without corrupting the running totals.
	sessionCost := sessionUsage.Cost
	if input.Cost.TotalCostUSD > 0 {
		sessionCost = input.Cost.TotalCostUSD
	}

	return themes.StatusData{
		ModelName:       formatModelName(input.Model.DisplayName),
		ModelType:       modelType,
		Version:         version,
		UpdateAvailable: updateAvailable,
		ProjectPath:     formatProjectPath(input.Workspace.CurrentDir),
		GitBranch:       gitInfo.Branch,
		GitStaged:       gitInfo.StagedCount,
		GitDirty:        gitInfo.DirtyCount,
		GitAhead:        gitInfo.AheadCount,
		GitBehind:       gitInfo.BehindCount,
		GitStash:        gitInfo.StashCount,
		GitSHA:          gitInfo.ShortSHA,
		TokenCount:      totalTokens,
		MessageCount:    sessionUsage.MessageCount,
		SessionTime:     totalHours,
		CacheHitRate:    cacheHitRate,
		TokensPerSec:    tokensPerSec,
		SessionCost:     sessionCost,
		DayCost:         dailyStats.TotalCost,
		MonthCost:       monthlyStats.TotalCost,
		WeekCost:        weeklyStats.TotalCost,
		BurnRate:        burnRate,
		ContextUsed:     contextUsed,
		ContextPercent:  contextPercent,
		API5hrPercent:   api5hrPercent,
		API5hrTimeLeft:  api5hrTimeLeft,
		API7dayPercent:  api7dayPercent,
		API7dayTimeLeft: api7dayTimeLeft,
		LinesAdded:      input.Cost.TotalLinesAdded,
		LinesRemoved:    input.Cost.TotalLinesRemoved,
	}
}

// resolveVersion returns the Claude Code version, preferring the value Claude
// Code now passes in the JSON input (no subprocess) and only shelling out to
// `claude --version` as a fallback for older Claude Code. The bool reports
// whether an update is available.
func resolveVersion(inputVersion string) (string, bool) {
	version := formatCCVersion(inputVersion)
	if version == "" {
		version = claudeCodeVersionFromCLI()
	}
	return version, isUpdateAvailable()
}

// formatCCVersion normalizes a Claude Code version string to a "vX.Y.Z" form.
// Returns "" when the input is empty so callers can fall back.
func formatCCVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return v
}

// claudeCodeVersionFromCLI shells out to `claude --version` (fallback path).
func claudeCodeVersionFromCLI() string {
	cmd := exec.Command("claude", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "v?.?.?"
	}
	version := strings.TrimSpace(string(output))
	version = strings.TrimPrefix(version, "claude ")
	version = strings.TrimSuffix(version, " (Claude Code)")
	return formatCCVersion(version)
}

// isUpdateAvailable reports whether Claude Code has flagged a pending update.
func isUpdateAvailable() bool {
	homeDir, _ := os.UserHomeDir()
	updateFile := filepath.Join(homeDir, ".claude", ".update_available")
	_, err := os.Stat(updateFile)
	return err == nil
}

// formatModelName formats model name
func formatModelName(displayName string) string {
	// 移除 "Claude " 前綴，只保留模型名稱和版本
	name := strings.TrimPrefix(displayName, "Claude ")
	if name != "" {
		return name
	}
	return displayName
}

// getModelType gets model type from a model's display name.
func getModelType(displayName string) string {
	for key := range modelPricing {
		if strings.Contains(displayName, key) {
			return key
		}
	}
	return "Sonnet"
}

// getModelTypeFromInput resolves the pricing family, preferring the stable
// model.id (e.g. "claude-opus-4-8") over the human display name. Falls back to
// display-name matching for older Claude Code that doesn't send model.id.
func getModelTypeFromInput(input Input) string {
	if id := strings.ToLower(input.Model.ID); id != "" {
		for key := range modelPricing {
			if strings.Contains(id, strings.ToLower(key)) {
				return key
			}
		}
	}
	return getModelType(input.Model.DisplayName)
}

// isFastMode reports whether the session uses Claude Code's Fast mode, which is
// billed at ~2x the standard per-token rates. Detected via the model id suffix
// (e.g. "claude-opus-4-8-fast") or a "fast" token / "(fast)" marker in the
// display name. Uses boundary matching rather than a bare substring so names
// like "Steadfast" or "fastpath" don't false-positive.
func isFastMode(input Input) bool {
	if strings.HasSuffix(strings.ToLower(input.Model.ID), "-fast") {
		return true
	}
	display := strings.ToLower(input.Model.DisplayName)
	if strings.Contains(display, "(fast)") {
		return true
	}
	return slices.Contains(strings.Fields(display), "fast")
}

// formatProjectPath formats project path
func formatProjectPath(fullPath string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Base(fullPath)
	}
	if strings.HasPrefix(fullPath, homeDir) {
		return "~" + fullPath[len(homeDir):]
	}
	return fullPath
}

// getOAuthToken gets OAuth token
func getOAuthToken() string {
	homeDir, _ := os.UserHomeDir()
	credFile := filepath.Join(homeDir, ".claude", ".credentials.json")
	output, err := os.ReadFile(credFile)

	if err != nil {
		cmd := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w")
		output, err = cmd.Output()
		if err != nil {
			return ""
		}
	}

	var creds struct {
		ClaudeAiOauth struct {
			AccessToken string `json:"accessToken"`
		} `json:"claudeAiOauth"`
	}
	if err := json.Unmarshal(output, &creds); err != nil {
		return ""
	}

	return creds.ClaudeAiOauth.AccessToken
}

// apiUsageCachePath returns the file path for the API usage cache.
func apiUsageCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".claude", "session-tracker", "api-usage-cache.json")
}

// apiUsageFromInput builds an *APIUsage from the rate_limits Claude Code now
// passes in the stdin JSON, avoiding a network round-trip. Returns nil when no
// window is present (older Claude Code, or non-subscriber sessions), in which
// case callers fall back to fetchAPIUsage. Each window is independently
// optional, signalled by a non-zero resets_at.
func apiUsageFromInput(input Input) *APIUsage {
	rl := input.RateLimits
	if rl.FiveHour.ResetsAt == 0 && rl.SevenDay.ResetsAt == 0 {
		return nil
	}

	var usage APIUsage
	usage.FiveHour.Utilization = rl.FiveHour.UsedPercentage
	if rl.FiveHour.ResetsAt > 0 {
		usage.FiveHour.ResetsAt = strconv.FormatInt(rl.FiveHour.ResetsAt, 10)
	}
	usage.SevenDay.Utilization = rl.SevenDay.UsedPercentage
	if rl.SevenDay.ResetsAt > 0 {
		usage.SevenDay.ResetsAt = strconv.FormatInt(rl.SevenDay.ResetsAt, 10)
	}
	return &usage
}

// fetchAPIUsage fetches API usage using the configured method.
// Dispatches to haiku probe or oauth usage endpoint based on config.usage_api.
// Results are cached to a file so that separate process invocations share the same cache.
func fetchAPIUsage() *APIUsage {
	cachePath := apiUsageCachePath()

	// Try file-based cache first
	if data, err := os.ReadFile(cachePath); err == nil {
		var cached APIUsageCache
		if err := json.Unmarshal(data, &cached); err == nil {
			if time.Since(cached.CachedAt) < 5*time.Minute {
				return &cached.Usage
			}
		}
	}

	token := getOAuthToken()
	if token == "" {
		return nil
	}

	config := loadConfig()
	var usage *APIUsage
	if config.UsageAPI == "haiku_probe" {
		usage = fetchViaHaikuProbe(token)
	} else {
		usage = fetchViaOAuthUsage(token)
	}

	if usage == nil {
		return nil
	}

	// Write to file cache
	cached := APIUsageCache{Usage: *usage, CachedAt: time.Now()}
	if data, err := json.Marshal(cached); err == nil {
		os.MkdirAll(filepath.Dir(cachePath), 0755)
		os.WriteFile(cachePath, data, 0644)
	}

	return usage
}

// fetchViaHaikuProbe sends a minimal Haiku request and reads rate limit headers.
func fetchViaHaikuProbe(token string) *APIUsage {
	client := &http.Client{Timeout: 5 * time.Second}
	body := `{"model":"claude-haiku-4-5-20251001","max_tokens":1,"messages":[{"role":"user","content":"hi"}]}`
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", strings.NewReader(body))
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	usage := APIUsage{}
	if v := resp.Header.Get("anthropic-ratelimit-unified-5h-utilization"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			usage.FiveHour.Utilization = f * 100 // header is 0.0-1.0, convert to percent
		}
	}
	if v := resp.Header.Get("anthropic-ratelimit-unified-5h-reset"); v != "" {
		usage.FiveHour.ResetsAt = v
	}
	if v := resp.Header.Get("anthropic-ratelimit-unified-7d-utilization"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			usage.SevenDay.Utilization = f * 100
		}
	}
	if v := resp.Header.Get("anthropic-ratelimit-unified-7d-reset"); v != "" {
		usage.SevenDay.ResetsAt = v
	}

	if usage.FiveHour.ResetsAt == "" && usage.SevenDay.ResetsAt == "" {
		return nil
	}
	return &usage
}

// fetchViaOAuthUsage calls the dedicated /api/oauth/usage endpoint.
func fetchViaOAuthUsage(token string) *APIUsage {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var usage APIUsage
	if err := json.Unmarshal(body, &usage); err != nil {
		return nil
	}

	if usage.FiveHour.ResetsAt == "" && usage.SevenDay.ResetsAt == "" {
		return nil
	}
	return &usage
}

// getGitInfo gets Git information
func getGitInfo() GitInfo {
	result := GitInfo{}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "rev-parse", "--git-dir")
		if err := cmd.Run(); err != nil {
			return result
		}
	}

	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return result
	}
	result.Branch = strings.TrimSpace(string(output))

	// --branch adds a leading "## <branch>...<upstream> [ahead N, behind M]"
	// header we parse for ahead/behind, at no extra git process cost.
	cmd = exec.Command("git", "status", "--porcelain", "--branch")
	output, err = cmd.Output()
	if err != nil {
		return result
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			result.AheadCount, result.BehindCount = parseAheadBehind(line)
			continue
		}
		if len(line) < 2 {
			continue
		}
		indexStatus := line[0]
		workTreeStatus := line[1]

		if indexStatus != ' ' && indexStatus != '?' {
			result.StagedCount++
		}
		if workTreeStatus != ' ' || indexStatus == '?' {
			result.DirtyCount++
		}
	}

	// Short HEAD SHA (empty on an unborn branch with no commits yet).
	if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		result.ShortSHA = strings.TrimSpace(string(out))
	}

	// Stash entry count (one line per stash; empty output means zero).
	if out, err := exec.Command("git", "stash", "list").Output(); err == nil {
		if s := strings.TrimSpace(string(out)); s != "" {
			result.StashCount = strings.Count(s, "\n") + 1
		}
	}

	return result
}

// parseAheadBehind extracts ahead/behind commit counts from a `git status
// --porcelain --branch` header line, e.g.
//
//	## main...origin/main [ahead 2, behind 1]
//
// Returns (0, 0) when the branch is level with upstream, has no upstream, or the
// bracket segment is absent/malformed.
func parseAheadBehind(header string) (ahead, behind int) {
	open := strings.LastIndex(header, "[")
	closeIdx := strings.LastIndex(header, "]")
	if open == -1 || closeIdx == -1 || closeIdx < open {
		return 0, 0
	}
	inside := header[open+1 : closeIdx] // "ahead 2, behind 1"
	for _, part := range strings.Split(inside, ",") {
		fields := strings.Fields(part)
		if len(fields) != 2 {
			continue
		}
		n, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		switch fields[0] {
		case "ahead":
			ahead = n
		case "behind":
			behind = n
		}
	}
	return ahead, behind
}

// updateSession updates session
func updateSession(sessionID string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	os.MkdirAll(sessionsDir, 0755)

	sessionFile := filepath.Join(sessionsDir, sessionID+".json")
	currentTime := time.Now().Unix()
	today := time.Now().Format("2006-01-02")

	var session Session

	if data, err := os.ReadFile(sessionFile); err == nil {
		json.Unmarshal(data, &session)
		if session.Date != today {
			session.Date = today
		}
	} else {
		session = Session{
			ID:            sessionID,
			Date:          today,
			Start:         currentTime,
			LastHeartbeat: currentTime,
			TotalSeconds:  0,
			Intervals:     []Interval{{Start: currentTime, End: nil}},
		}
	}

	gap := currentTime - session.LastHeartbeat
	session.LastHeartbeat = currentTime

	if gap < 600 {
		if len(session.Intervals) > 0 {
			session.Intervals[len(session.Intervals)-1].End = &currentTime
		}
	} else {
		session.Intervals = append(session.Intervals, Interval{
			Start: currentTime,
			End:   &currentTime,
		})
	}

	var total int64
	for _, interval := range session.Intervals {
		if interval.End != nil {
			total += *interval.End - interval.Start
		}
	}
	session.TotalSeconds = total

	if data, err := json.Marshal(session); err == nil {
		os.WriteFile(sessionFile, data, 0644)
	}
}

// calculateTotalHours calculates total hours
func calculateTotalHours(currentSessionID string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "0m"
	}

	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return "0m"
	}

	var totalSeconds int64
	today := time.Now().Format("2006-01-02")

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		sessionFile := filepath.Join(sessionsDir, entry.Name())
		data, err := os.ReadFile(sessionFile)
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		if session.Date == today {
			totalSeconds += session.TotalSeconds
		}
	}

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60

	if hours > 0 {
		return fmt.Sprintf("%dh%02dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// calculateSessionUsage calculates session usage
func calculateSessionUsage(transcriptPath, sessionID, modelType string, fast bool) SessionUsageResult {
	result := SessionUsageResult{}

	if transcriptPath == "" {
		return result
	}

	file, err := os.Open(transcriptPath)
	if err != nil {
		return result
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxScanTokenSize = 1024 * 1024
	buf := make([]byte, 0, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	var sessionStart time.Time
	var lastTime time.Time

	// Claude Code streams each assistant message to the transcript multiple
	// times (partial + final writes) under the same message.id with identical
	// final usage numbers. Summing every line double-counts tokens (~2x+), which
	// inflates the displayed token count and cache-hit rate. Dedup by message.id,
	// keeping the last occurrence; entries without an id (rare) can't be deduped
	// and are summed directly. See calculateSessionUsage dedup test.
	type usageTokens struct {
		input, output, cacheRead, cacheWrite int64
	}
	byMsgID := make(map[string]usageTokens)
	var order []string // preserve first-seen order for deterministic summing

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		if sid, ok := data["sessionId"].(string); !ok || sid != sessionID {
			continue
		}

		if isSidechain, ok := data["isSidechain"].(bool); ok && isSidechain {
			continue
		}

		if ts, ok := data["timestamp"].(string); ok {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				if sessionStart.IsZero() {
					sessionStart = t
				}
				lastTime = t
			}
		}

		if msgType, ok := data["type"].(string); ok && msgType == "user" {
			result.MessageCount++
		}

		if message, ok := data["message"].(map[string]interface{}); ok {
			if usage, ok := message["usage"].(map[string]interface{}); ok {
				var u usageTokens
				if input, ok := usage["input_tokens"].(float64); ok {
					u.input = int64(input)
				}
				if output, ok := usage["output_tokens"].(float64); ok {
					u.output = int64(output)
				}
				if cacheRead, ok := usage["cache_read_input_tokens"].(float64); ok {
					u.cacheRead = int64(cacheRead)
				}
				if cacheCreation, ok := usage["cache_creation_input_tokens"].(float64); ok {
					u.cacheWrite = int64(cacheCreation)
				}

				if id, ok := message["id"].(string); ok && id != "" {
					if _, seen := byMsgID[id]; !seen {
						order = append(order, id)
					}
					byMsgID[id] = u // last write wins
				} else {
					// No message id to dedup on; count it once, here.
					result.InputTokens += u.input
					result.OutputTokens += u.output
					result.CacheReadTokens += u.cacheRead
					result.CacheWriteTokens += u.cacheWrite
				}
			}
		}
	}

	for _, id := range order {
		u := byMsgID[id]
		result.InputTokens += u.input
		result.OutputTokens += u.output
		result.CacheReadTokens += u.cacheRead
		result.CacheWriteTokens += u.cacheWrite
	}

	if !sessionStart.IsZero() && !lastTime.IsZero() {
		result.Duration = lastTime.Sub(sessionStart)
	}

	result.Cost = calculateCostMode(result, modelType, fast)

	return result
}

// calculateCost calculates cost at standard (non-Fast) rates.
func calculateCost(usage SessionUsageResult, modelType string) float64 {
	return calculateCostMode(usage, modelType, false)
}

// calculateCostMode calculates cost from token counts. When fast is true the
// per-token rates are doubled to approximate Claude Code's Fast mode (~2x
// standard pricing). This is a fallback estimate only — when Claude Code sends
// cost.total_cost_usd we use that authoritative value instead.
func calculateCostMode(usage SessionUsageResult, modelType string, fast bool) float64 {
	pricing, ok := modelPricing[modelType]
	if !ok {
		pricing = modelPricing["Sonnet"]
	}

	mult := 1.0
	if fast {
		mult = 2.0
	}

	cost := float64(usage.InputTokens) * pricing.Input * mult / 1000000
	cost += float64(usage.OutputTokens) * pricing.Output * mult / 1000000
	cost += float64(usage.CacheReadTokens) * pricing.CacheRead * mult / 1000000
	cost += float64(usage.CacheWriteTokens) * pricing.CacheWrite * mult / 1000000

	return cost
}

// applyCostDelta accumulates a session's monotonically-increasing cumulative
// cost into a stats bucket. It records the last-known cumulative cost per
// session and adds only the positive delta, so repeated renders within a
// session never double-count, and switching the cost source (transcript
// estimate ↔ Claude Code's cost.total_cost_usd) cannot corrupt the running
// total. A lower incoming value (e.g. the authoritative value landing below a
// prior transcript over-estimate) is conservatively ignored — totals stay
// correct on the next rise rather than risking a double-count.
func applyCostDelta(stats *UsageStats, sessionID string, sessionCost float64) {
	if stats.SessionCosts == nil {
		stats.SessionCosts = make(map[string]float64)
	}
	delta := sessionCost - stats.SessionCosts[sessionID]
	if delta > 0 {
		stats.TotalCost += delta
		stats.SessionCosts[sessionID] = sessionCost
	}
}

// getDailyStats gets daily stats
func getDailyStats() UsageStats {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")
	today := time.Now().Format("2006-01-02")
	statsFile := filepath.Join(statsDir, "daily-"+today+".json")

	var stats UsageStats
	if data, err := os.ReadFile(statsFile); err == nil {
		json.Unmarshal(data, &stats)
	}
	stats.Date = today

	return stats
}

// getWeeklyStats gets weekly stats
func getWeeklyStats() UsageStats {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")

	statsFile := filepath.Join(statsDir, "weekly-"+weekStart+".json")

	var stats UsageStats
	if data, err := os.ReadFile(statsFile); err == nil {
		json.Unmarshal(data, &stats)
	}
	stats.Week = weekStart

	return stats
}

// getMonthlyStats gets monthly stats
func getMonthlyStats() UsageStats {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")

	monthKey := time.Now().Format("2006-01")
	statsFile := filepath.Join(statsDir, "monthly-"+monthKey+".json")

	var stats UsageStats
	if data, err := os.ReadFile(statsFile); err == nil {
		json.Unmarshal(data, &stats)
	}

	return stats
}

// updateDailyStats updates daily stats
func updateDailyStats(sessionID string, data themes.StatusData, modelType string) {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")
	os.MkdirAll(statsDir, 0755)

	today := time.Now().Format("2006-01-02")
	dailyFile := filepath.Join(statsDir, "daily-"+today+".json")

	var dailyStats UsageStats
	if fileData, err := os.ReadFile(dailyFile); err == nil {
		json.Unmarshal(fileData, &dailyStats)
	}

	applyCostDelta(&dailyStats, sessionID, data.SessionCost)

	dailyStats.Date = today
	dailyStats.LastUpdated = time.Now().Unix()

	if fileData, err := json.Marshal(dailyStats); err == nil {
		os.WriteFile(dailyFile, fileData, 0644)
	}

	updateWeeklyStats(sessionID, data.SessionCost)
	updateMonthlyStats(sessionID, data.SessionCost)
}

// updateWeeklyStats updates weekly stats
func updateWeeklyStats(sessionID string, sessionCost float64) {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := now.AddDate(0, 0, -(weekday - 1)).Format("2006-01-02")

	weeklyFile := filepath.Join(statsDir, "weekly-"+weekStart+".json")

	var weeklyStats UsageStats
	if data, err := os.ReadFile(weeklyFile); err == nil {
		json.Unmarshal(data, &weeklyStats)
	}

	applyCostDelta(&weeklyStats, sessionID, sessionCost)

	weeklyStats.Week = weekStart
	weeklyStats.LastUpdated = time.Now().Unix()

	if data, err := json.Marshal(weeklyStats); err == nil {
		os.WriteFile(weeklyFile, data, 0644)
	}
}

// updateMonthlyStats updates monthly stats
func updateMonthlyStats(sessionID string, sessionCost float64) {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")

	monthKey := time.Now().Format("2006-01")
	monthlyFile := filepath.Join(statsDir, "monthly-"+monthKey+".json")

	var monthlyStats UsageStats
	if data, err := os.ReadFile(monthlyFile); err == nil {
		json.Unmarshal(data, &monthlyStats)
	}

	applyCostDelta(&monthlyStats, sessionID, sessionCost)

	monthlyStats.LastUpdated = time.Now().Unix()

	if data, err := json.Marshal(monthlyStats); err == nil {
		os.WriteFile(monthlyFile, data, 0644)
	}
}

// calculateBurnRateValue calculates burn rate value
func calculateBurnRateValue(dailyStats UsageStats) float64 {
	homeDir, _ := os.UserHomeDir()
	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	entries, _ := os.ReadDir(sessionsDir)

	var totalSeconds int64
	today := time.Now().Format("2006-01-02")

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		sessionFile := filepath.Join(sessionsDir, entry.Name())
		data, _ := os.ReadFile(sessionFile)

		var session Session
		if err := json.Unmarshal(data, &session); err == nil && session.Date == today {
			totalSeconds += session.TotalSeconds
		}
	}

	if totalSeconds < 300 {
		return 0
	}

	hours := float64(totalSeconds) / 3600
	return dailyStats.TotalCost / hours
}

// formatTimeLeftShort formats time left in short form
func formatTimeLeftShort(timeStr string) string {
	var t time.Time
	// Try Unix epoch seconds first (from rate limit headers), then ISO 8601 (from usage API)
	if epoch, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		t = time.Unix(epoch, 0)
	} else if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
		t = parsed
	} else {
		return "?"
	}

	now := time.Now()
	diff := t.Sub(now)

	if diff <= 0 {
		return "now"
	}

	days := int(diff.Hours() / 24)
	hours := int(diff.Hours()) % 24
	minutes := int(diff.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
