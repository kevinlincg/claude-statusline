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
	"sort"
	"strings"
	"sync"
	"time"

	"statusline/themes"
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
		Input:      5.0,   // Opus 4.5: $5 per 1M input tokens
		Output:     25.0,  // Opus 4.5: $25 per 1M output tokens
		CacheRead:  0.5,   // Opus 4.5: $0.50 per 1M cache read tokens
		CacheWrite: 6.25,  // Opus 4.5: $6.25 per 1M cache write tokens (5m)
	},
	"Sonnet": {
		Input:      3.0,   // Sonnet 4/4.5: $3 per 1M input tokens
		Output:     15.0,  // Sonnet 4/4.5: $15 per 1M output tokens
		CacheRead:  0.3,   // Sonnet 4/4.5: $0.30 per 1M cache read tokens
		CacheWrite: 3.75,  // Sonnet 4/4.5: $3.75 per 1M cache write tokens (5m)
	},
	"Haiku": {
		Input:      1.0,   // Haiku 4.5: $1 per 1M input tokens
		Output:     5.0,   // Haiku 4.5: $5 per 1M output tokens
		CacheRead:  0.1,   // Haiku 4.5: $0.10 per 1M cache read tokens
		CacheWrite: 1.25,  // Haiku 4.5: $1.25 per 1M cache write tokens (5m)
	},
}

// Input data structure
type Input struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	SessionID string `json:"session_id"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	TranscriptPath string `json:"transcript_path,omitempty"`
}

// Config structure
type Config struct {
	Theme string `json:"theme"`
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

// Cache
var (
	apiUsageCache   *APIUsage
	apiUsageExpires time.Time
	cacheMutex      sync.RWMutex
)

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

	// Get model type
	modelType := getModelType(input.Model.DisplayName)

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
		ModelName:       "Opus 4.5",
		ModelType:       "Opus",
		Version:         "v1.0.75",
		UpdateAvailable: true,
		ProjectPath:     "~/cookys/project",
		GitBranch:       "main",
		GitStaged:       3,
		GitDirty:        5,
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

		if n == 1 {
			switch buf[0] {
			case 'q', 'Q', 27: // q or Escape
				fmt.Print("\033[2J\033[H")
				fmt.Print("Canceled\r\n")
				return
			case 13, 10: // Enter
				fmt.Print("\033[2J\033[H")
				saveThemeConfig(themeList[selectedIndex].Name())
				fmt.Printf("Theme set to: %s\r\n", themeList[selectedIndex].Name())
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
func makeRaw(fd uintptr) ([]byte, error) {
	// Use stty to set raw mode
	cmd := exec.Command("stty", "-F", "/dev/stdin", "raw", "-echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		// macOS uses different syntax
		cmd = exec.Command("stty", "raw", "-echo")
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
	return nil, nil
}

func restore(fd uintptr, oldState []byte) {
	// Restore terminal settings
	cmd := exec.Command("stty", "-F", "/dev/stdin", "sane")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		// macOS
		cmd = exec.Command("stty", "sane")
		cmd.Stdin = os.Stdin
		cmd.Run()
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
		ModelName:       "Opus 4.5",
		ModelType:       "Opus",
		Version:         "v1.0.75",
		UpdateAvailable: true,
		ProjectPath:     "~/cookys/project",
		GitBranch:       "main",
		GitStaged:       3,
		GitDirty:        5,
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

	homeDir, _ := os.UserHomeDir()
	configFile := filepath.Join(homeDir, ".claude", "statusline-go", "config.json")

	config := Config{Theme: themeName}
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0644)

	fmt.Printf("Theme set to: %s\n", themeName)
}

// loadThemeConfig loads theme configuration
func loadThemeConfig() string {
	homeDir, _ := os.UserHomeDir()
	configFile := filepath.Join(homeDir, ".claude", "statusline-go", "config.json")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return "classic_framed" // Default theme
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return "classic_framed"
	}

	if config.Theme == "" {
		return "classic_framed"
	}

	return config.Theme
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

	go func() {
		defer wg.Done()
		sessionUsage := calculateSessionUsage(input.TranscriptPath, input.SessionID, modelType)
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
		apiUsage := fetchAPIUsage()
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

	// Calculate Context
	contextUsed := 0
	contextPercent := 0
	if input.TranscriptPath != "" {
		contextUsed = calculateContextUsage(input.TranscriptPath)
		contextPercent = int(float64(contextUsed) * 100.0 / 200000.0)
		if contextPercent > 100 {
			contextPercent = 100
		}
	}

	// Get monthly stats
	monthlyStats := getMonthlyStats()

	// Calculate burn rate
	burnRate := calculateBurnRateValue(dailyStats)

	// Get version and update status
	version, updateAvailable := getVersionInfo()

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

	return themes.StatusData{
		ModelName:       formatModelName(input.Model.DisplayName),
		ModelType:       modelType,
		Version:         version,
		UpdateAvailable: updateAvailable,
		ProjectPath:     formatProjectPath(input.Workspace.CurrentDir),
		GitBranch:       gitInfo.Branch,
		GitStaged:       gitInfo.StagedCount,
		GitDirty:        gitInfo.DirtyCount,
		TokenCount:      sessionUsage.InputTokens + sessionUsage.OutputTokens + sessionUsage.CacheReadTokens + sessionUsage.CacheWriteTokens,
		MessageCount:    sessionUsage.MessageCount,
		SessionTime:     totalHours,
		CacheHitRate:    cacheHitRate,
		SessionCost:     sessionUsage.Cost,
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
	}
}

// getVersionInfo gets version information
func getVersionInfo() (string, bool) {
	// Try to execute claude --version
	cmd := exec.Command("claude", "--version")
	output, err := cmd.Output()
	version := "v?.?.?"
	if err == nil {
		version = strings.TrimSpace(string(output))
		// Remove extra prefix and suffix
		version = strings.TrimPrefix(version, "claude ")
		version = strings.TrimSuffix(version, " (Claude Code)")
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
	}

	// Check if update available (check if file exists)
	homeDir, _ := os.UserHomeDir()
	updateFile := filepath.Join(homeDir, ".claude", ".update_available")
	_, updateAvailable := os.Stat(updateFile)

	return version, updateAvailable == nil
}

// formatModelName formats model name
func formatModelName(displayName string) string {
	if strings.Contains(displayName, "Opus") {
		return "Opus 4.5"
	} else if strings.Contains(displayName, "Sonnet") {
		return "Sonnet 4"
	} else if strings.Contains(displayName, "Haiku") {
		return "Haiku 3.5"
	}
	return displayName
}

// getModelType gets model type
func getModelType(displayName string) string {
	for key := range modelPricing {
		if strings.Contains(displayName, key) {
			return key
		}
	}
	return "Sonnet"
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

// fetchAPIUsage fetches API usage
func fetchAPIUsage() *APIUsage {
	cacheMutex.RLock()
	if apiUsageCache != nil && time.Now().Before(apiUsageExpires) {
		result := apiUsageCache
		cacheMutex.RUnlock()
		return result
	}
	cacheMutex.RUnlock()

	token := getOAuthToken()
	if token == "" {
		return nil
	}

	client := &http.Client{Timeout: 3 * time.Second}
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

	cacheMutex.Lock()
	apiUsageCache = &usage
	apiUsageExpires = time.Now().Add(30 * time.Second)
	cacheMutex.Unlock()

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

	cmd = exec.Command("git", "status", "--porcelain")
	output, err = cmd.Output()
	if err != nil {
		return result
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
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

	return result
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
func calculateSessionUsage(transcriptPath, sessionID, modelType string) SessionUsageResult {
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
				if input, ok := usage["input_tokens"].(float64); ok {
					result.InputTokens += int64(input)
				}
				if output, ok := usage["output_tokens"].(float64); ok {
					result.OutputTokens += int64(output)
				}
				if cacheRead, ok := usage["cache_read_input_tokens"].(float64); ok {
					result.CacheReadTokens += int64(cacheRead)
				}
				if cacheCreation, ok := usage["cache_creation_input_tokens"].(float64); ok {
					result.CacheWriteTokens += int64(cacheCreation)
				}
			}
		}
	}

	if !sessionStart.IsZero() && !lastTime.IsZero() {
		result.Duration = lastTime.Sub(sessionStart)
	}

	result.Cost = calculateCost(result, modelType)

	return result
}

// calculateCost calculates cost
func calculateCost(usage SessionUsageResult, modelType string) float64 {
	pricing, ok := modelPricing[modelType]
	if !ok {
		pricing = modelPricing["Sonnet"]
	}

	cost := float64(usage.InputTokens) * pricing.Input / 1000000
	cost += float64(usage.OutputTokens) * pricing.Output / 1000000
	cost += float64(usage.CacheReadTokens) * pricing.CacheRead / 1000000
	cost += float64(usage.CacheWriteTokens) * pricing.CacheWrite / 1000000

	return cost
}

// calculateContextUsage calculates context usage
func calculateContextUsage(transcriptPath string) int {
	file, err := os.Open(transcriptPath)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	const maxScanTokenSize = 1024 * 1024
	buf := make([]byte, 0, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	allLines := make([]string, 0)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	start := len(allLines) - 100
	if start < 0 {
		start = 0
	}
	lines := allLines[start:]

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		if sidechain, ok := data["isSidechain"]; ok {
			if isSide, ok := sidechain.(bool); ok && isSide {
				continue
			}
		}

		if message, ok := data["message"].(map[string]interface{}); ok {
			if usage, ok := message["usage"].(map[string]interface{}); ok {
				var total float64

				if input, ok := usage["input_tokens"].(float64); ok {
					total += input
				}
				if cacheRead, ok := usage["cache_read_input_tokens"].(float64); ok {
					total += cacheRead
				}
				if cacheCreation, ok := usage["cache_creation_input_tokens"].(float64); ok {
					total += cacheCreation
				}

				if total > 0 {
					return int(total)
				}
			}
		}
	}

	return 0
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

	if dailyStats.SessionCosts == nil {
		dailyStats.SessionCosts = make(map[string]float64)
	}

	lastKnownCost := dailyStats.SessionCosts[sessionID]
	delta := data.SessionCost - lastKnownCost
	if delta > 0 {
		dailyStats.TotalCost += delta
		dailyStats.SessionCosts[sessionID] = data.SessionCost
	}

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

	if weeklyStats.SessionCosts == nil {
		weeklyStats.SessionCosts = make(map[string]float64)
	}

	lastKnownCost := weeklyStats.SessionCosts[sessionID]
	delta := sessionCost - lastKnownCost
	if delta > 0 {
		weeklyStats.TotalCost += delta
		weeklyStats.SessionCosts[sessionID] = sessionCost
	}

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

	if monthlyStats.SessionCosts == nil {
		monthlyStats.SessionCosts = make(map[string]float64)
	}

	lastKnownCost := monthlyStats.SessionCosts[sessionID]
	delta := sessionCost - lastKnownCost
	if delta > 0 {
		monthlyStats.TotalCost += delta
		monthlyStats.SessionCosts[sessionID] = sessionCost
	}

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
func formatTimeLeftShort(isoTime string) string {
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
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
