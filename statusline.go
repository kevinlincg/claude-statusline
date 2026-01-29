package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ANSI 顏色定義
const (
	ColorReset  = "\033[0m"
	ColorGold   = "\033[38;2;195;158;83m"
	ColorCyan   = "\033[38;2;118;170;185m"
	ColorPink   = "\033[38;2;255;182;193m"
	ColorGreen  = "\033[38;2;152;195;121m"
	ColorGray   = "\033[38;2;64;64;64m"
	ColorSilver = "\033[38;2;192;192;192m"
	ColorOrange = "\033[38;2;255;165;0m"
	ColorPurple = "\033[38;2;186;133;217m"
	ColorBlue   = "\033[38;2;100;149;237m"
	ColorRed    = "\033[38;2;220;88;88m"
	ColorDim    = "\033[38;2;128;128;128m"
	ColorYellow = "\033[38;2;255;215;0m"

	ColorCtxGreen = "\033[38;2;108;167;108m"
	ColorCtxGold  = "\033[38;2;188;155;83m"
	ColorCtxRed   = "\033[38;2;185;102;82m"
)

// 模型價格 (per 1M tokens) - 2026 年 1 月更新
// 官方定價: https://platform.claude.com/docs/en/about-claude/pricing
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

// 模型圖示和顏色
var modelConfig = map[string][2]string{
	"Opus":   {ColorGold, "💛"},
	"Sonnet": {ColorCyan, "💠"},
	"Haiku":  {ColorPink, "🌸"},
}

// 輸入資料結構
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

// Session 資料結構
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

// Usage 統計結構
type UsageStats struct {
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	CacheReadTokens  int64   `json:"cache_read_tokens"`
	CacheWriteTokens int64   `json:"cache_write_tokens"`
	TotalCost        float64 `json:"total_cost"`
	MessageCount     int     `json:"message_count"`
	SessionCount     int     `json:"session_count"`
	Date             string  `json:"date"`
	Week             string  `json:"week"`
	LastUpdated      int64   `json:"last_updated"`
}

// API Usage 結構
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

// 結果通道資料
type Result struct {
	Type string
	Data interface{}
}

// GitInfo 包含 Git 狀態資訊
type GitInfo struct {
	Branch      string
	DirtyCount  int
	StagedCount int
}

// SessionUsageResult 包含 session 的用量資訊
type SessionUsageResult struct {
	InputTokens      int64
	OutputTokens     int64
	CacheReadTokens  int64
	CacheWriteTokens int64
	Cost             float64
	MessageCount     int
	Duration         time.Duration
}

// 簡單快取
var (
	gitBranchCache   string
	gitBranchExpires time.Time
	apiUsageCache    *APIUsage
	apiUsageExpires  time.Time
	cacheMutex       sync.RWMutex
)

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode input: %v\n", err)
		os.Exit(1)
	}

	// 取得模型類型
	modelType := getModelType(input.Model.DisplayName)

	// 建立結果通道
	results := make(chan Result, 10)
	var wg sync.WaitGroup

	// 並行獲取各種資訊
	wg.Add(7)

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
		contextInfo := analyzeContext(input.TranscriptPath)
		results <- Result{"context", contextInfo}
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

	// 等待所有 goroutines 完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集結果
	var (
		gitInfo      GitInfo
		totalHours   string
		contextUsage string
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
		case "context":
			contextUsage = result.Data.(string)
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

	// 更新 session 和統計
	updateSession(input.SessionID)
	updateDailyStats(sessionUsage, modelType)

	// 格式化輸出
	modelDisplay := formatModel(input.Model.DisplayName)
	projectName := filepath.Base(input.Workspace.CurrentDir)
	gitDisplay := formatGitInfo(gitInfo)

	// 第一行：基本資訊
	fmt.Printf("%s[%s] 📂 %s%s%s | %s%s\n",
		ColorReset, modelDisplay, projectName, gitDisplay,
		contextUsage, totalHours, ColorReset)

	// 第二行：API 用量限制 (Session + Weekly)
	apiUsageInfo := formatAPIUsage(apiUsage)
	fmt.Printf("%s│ %s%s\n", ColorDim, apiUsageInfo, ColorReset)

	// 第三行：本地統計 (Session tokens/cost + 燒錢速度 + 今日/週成本 + Cache 命中率)
	sessionInfo := formatSessionUsage(sessionUsage)
	burnRate := calculateBurnRate(dailyStats)
	costInfo := formatCostStats(dailyStats, weeklyStats)
	cacheHitRate := formatCacheHitRate(sessionUsage)
	fmt.Printf("%s│ %s | %s | %s | %s%s\n", ColorDim, sessionInfo, burnRate, costInfo, cacheHitRate, ColorReset)
}

// 獲取 OAuth Token (支援 Linux 和 macOS)
func getOAuthToken() string {
	var output []byte
	var err error

	// 先嘗試 Linux: 從 ~/.claude/.credentials.json 讀取
	homeDir, _ := os.UserHomeDir()
	credFile := filepath.Join(homeDir, ".claude", ".credentials.json")
	output, err = os.ReadFile(credFile)

	// 如果檔案不存在，嘗試 macOS Keychain
	if err != nil {
		cmd := exec.Command("security", "find-generic-password", "-s", "Claude Code-credentials", "-w")
		output, err = cmd.Output()
		if err != nil {
			return ""
		}
	}

	// 解析 JSON 取得 access_token (nested structure)
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

// 獲取 API Usage
func fetchAPIUsage() *APIUsage {
	// 檢查快取
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

	// 更新快取 (30秒)
	cacheMutex.Lock()
	apiUsageCache = &usage
	apiUsageExpires = time.Now().Add(30 * time.Second)
	cacheMutex.Unlock()

	return &usage
}

// 格式化 API Usage
func formatAPIUsage(usage *APIUsage) string {
	if usage == nil {
		return fmt.Sprintf("%s⚠️  API usage unavailable%s", ColorDim, ColorReset)
	}

	// Session (5-hour) 用量
	sessionPct := int(usage.FiveHour.Utilization)
	sessionBar := generateUsageBar(sessionPct, 8)
	sessionReset := formatResetTime(usage.FiveHour.ResetsAt)
	sessionColor := getUsageColor(sessionPct)

	// Weekly (7-day) 用量
	weeklyPct := int(usage.SevenDay.Utilization)
	weeklyBar := generateUsageBar(weeklyPct, 8)
	weeklyReset := formatResetTime(usage.SevenDay.ResetsAt)
	weeklyColor := getUsageColor(weeklyPct)

	return fmt.Sprintf("⏱️ Session %s %s%d%%%s ↻%s | 📅 Week %s %s%d%%%s ↻%s",
		sessionBar, sessionColor, sessionPct, ColorReset, sessionReset,
		weeklyBar, weeklyColor, weeklyPct, ColorReset, weeklyReset)
}

// 生成用量進度條
func generateUsageBar(percentage, width int) string {
	filled := percentage * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled
	color := getUsageColor(percentage)

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(ColorReset)
	}
	if empty > 0 {
		bar.WriteString(ColorGray)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(ColorReset)
	}

	return bar.String()
}

// 獲取用量顏色
func getUsageColor(percentage int) string {
	if percentage < 50 {
		return ColorGreen
	} else if percentage < 75 {
		return ColorYellow
	} else if percentage < 90 {
		return ColorOrange
	}
	return ColorRed
}

// 格式化 Reset 時間
func formatResetTime(isoTime string) string {
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return "?"
	}

	// 轉換為本地時間
	local := t.Local()
	now := time.Now()

	// 如果是今天
	if local.Day() == now.Day() && local.Month() == now.Month() {
		return local.Format("3:04pm")
	}

	// 如果是明天
	tomorrow := now.AddDate(0, 0, 1)
	if local.Day() == tomorrow.Day() && local.Month() == tomorrow.Month() {
		return "明天" + local.Format("3:04pm")
	}

	// 其他日期
	return local.Format("1/2 3pm")
}

// 獲取模型類型
func getModelType(displayName string) string {
	for key := range modelPricing {
		if strings.Contains(displayName, key) {
			return key
		}
	}
	return "Sonnet" // 預設
}

// 格式化模型顯示
func formatModel(model string) string {
	for key, config := range modelConfig {
		if strings.Contains(model, key) {
			color := config[0]
			icon := config[1]
			return fmt.Sprintf("%s%s %s%s", color, icon, model, ColorReset)
		}
	}
	return model
}

// 獲取 Git 資訊（分支名稱 + 狀態）
func getGitInfo() GitInfo {
	result := GitInfo{}

	// 檢查是否在 Git 倉庫中
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "rev-parse", "--git-dir")
		if err := cmd.Run(); err != nil {
			return result
		}
	}

	// 獲取分支名稱
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return result
	}
	result.Branch = strings.TrimSpace(string(output))

	// 獲取未暫存的修改數量 (modified, deleted, untracked)
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

		// 已暫存的檔案 (index 有狀態)
		if indexStatus != ' ' && indexStatus != '?' {
			result.StagedCount++
		}
		// 未暫存的修改 (工作區有狀態或是 untracked)
		if workTreeStatus != ' ' || indexStatus == '?' {
			result.DirtyCount++
		}
	}

	return result
}

// 格式化 Git 資訊
func formatGitInfo(info GitInfo) string {
	if info.Branch == "" {
		return ""
	}

	result := fmt.Sprintf(" %s⚡ %s%s", ColorCyan, info.Branch, ColorReset)

	// 顯示 Git 狀態
	if info.StagedCount > 0 || info.DirtyCount > 0 {
		statusStr := ""
		if info.StagedCount > 0 {
			statusStr += fmt.Sprintf("%s+%d%s", ColorGreen, info.StagedCount, ColorReset)
		}
		if info.DirtyCount > 0 {
			if statusStr != "" {
				statusStr += "/"
			}
			statusStr += fmt.Sprintf("%s~%d%s", ColorOrange, info.DirtyCount, ColorReset)
		}
		result += fmt.Sprintf(" [%s]", statusStr)
	}

	return result
}

// 更新 Session
func updateSession(sessionID string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return
	}

	sessionFile := filepath.Join(sessionsDir, sessionID+".json")
	currentTime := time.Now().Unix()
	today := time.Now().Format("2006-01-02")

	var session Session

	if data, err := os.ReadFile(sessionFile); err == nil {
		json.Unmarshal(data, &session)
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

// 計算總時數
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
	activeSessions := 0
	today := time.Now().Format("2006-01-02")
	currentTime := time.Now().Unix()

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

			if currentTime-session.LastHeartbeat < 600 {
				activeSessions++
			}
		}
	}

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60

	var timeStr string
	if hours > 0 {
		timeStr = fmt.Sprintf("%dh", hours)
		if minutes > 0 {
			timeStr += fmt.Sprintf("%dm", minutes)
		}
	} else {
		timeStr = fmt.Sprintf("%dm", minutes)
	}

	if activeSessions > 1 {
		return fmt.Sprintf("%s [%d sessions]", timeStr, activeSessions)
	}
	return timeStr
}

// 計算 Session 用量
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

		// 檢查是否為當前 session
		if sid, ok := data["sessionId"].(string); !ok || sid != sessionID {
			continue
		}

		// 跳過 sidechain
		if isSidechain, ok := data["isSidechain"].(bool); ok && isSidechain {
			continue
		}

		// 解析時間戳
		if ts, ok := data["timestamp"].(string); ok {
			if t, err := time.Parse(time.RFC3339, ts); err == nil {
				if sessionStart.IsZero() {
					sessionStart = t
				}
				lastTime = t
			}
		}

		// 統計訊息數
		if msgType, ok := data["type"].(string); ok && msgType == "user" {
			result.MessageCount++
		}

		// 提取 usage 資料
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

	// 計算持續時間
	if !sessionStart.IsZero() && !lastTime.IsZero() {
		result.Duration = lastTime.Sub(sessionStart)
	}

	// 計算成本
	result.Cost = calculateCost(result, modelType)

	return result
}

// 計算成本
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

// 格式化 Session 用量
func formatSessionUsage(usage SessionUsageResult) string {
	totalTokens := usage.InputTokens + usage.OutputTokens + usage.CacheReadTokens + usage.CacheWriteTokens

	tokenStr := formatTokenCount(totalTokens)
	costStr := formatCost(usage.Cost)

	durationStr := ""
	if usage.Duration > 0 {
		if usage.Duration.Hours() >= 1 {
			durationStr = fmt.Sprintf(" ⏱️%dh%dm", int(usage.Duration.Hours()), int(usage.Duration.Minutes())%60)
		} else {
			durationStr = fmt.Sprintf(" ⏱️%dm", int(usage.Duration.Minutes()))
		}
	}

	msgStr := ""
	if usage.MessageCount > 0 {
		msgStr = fmt.Sprintf(" 💬%d", usage.MessageCount)
	}

	return fmt.Sprintf("%s🔤%s%s %s💰%s%s%s%s",
		ColorPurple, tokenStr, ColorReset,
		ColorGreen, costStr, ColorReset,
		durationStr, msgStr)
}

// 格式化 Token 數量
func formatTokenCount(tokens int64) string {
	if tokens >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
	} else if tokens >= 1000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}

// 格式化成本
func formatCost(cost float64) string {
	if cost >= 1.0 {
		return fmt.Sprintf("$%.2f", cost)
	} else if cost >= 0.01 {
		return fmt.Sprintf("$%.3f", cost)
	}
	return fmt.Sprintf("$%.4f", cost)
}

// 獲取每日統計
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

// 獲取每週統計
func getWeeklyStats() UsageStats {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")

	// 計算當週開始日期 (週一)
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

// 更新每日統計
func updateDailyStats(sessionUsage SessionUsageResult, modelType string) {
	homeDir, _ := os.UserHomeDir()
	statsDir := filepath.Join(homeDir, ".claude", "session-tracker", "stats")
	os.MkdirAll(statsDir, 0755)

	today := time.Now().Format("2006-01-02")
	dailyFile := filepath.Join(statsDir, "daily-"+today+".json")

	// 讀取現有統計
	var dailyStats UsageStats
	if data, err := os.ReadFile(dailyFile); err == nil {
		json.Unmarshal(data, &dailyStats)
	}

	// 更新統計
	dailyStats.Date = today
	dailyStats.TotalCost += sessionUsage.Cost * 0.01
	dailyStats.LastUpdated = time.Now().Unix()

	// 儲存
	if data, err := json.Marshal(dailyStats); err == nil {
		os.WriteFile(dailyFile, data, 0644)
	}

	// 同時更新每週統計
	updateWeeklyStats(sessionUsage, modelType)
}

// 更新每週統計
func updateWeeklyStats(sessionUsage SessionUsageResult, modelType string) {
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

	weeklyStats.Week = weekStart
	weeklyStats.TotalCost += sessionUsage.Cost * 0.01
	weeklyStats.LastUpdated = time.Now().Unix()

	if data, err := json.Marshal(weeklyStats); err == nil {
		os.WriteFile(weeklyFile, data, 0644)
	}
}

// 計算燒錢速度
func calculateBurnRate(dailyStats UsageStats) string {
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

	if totalSeconds < 60 {
		return fmt.Sprintf("%s🔥--/hr%s", ColorRed, ColorReset)
	}

	hours := float64(totalSeconds) / 3600
	rate := dailyStats.TotalCost / hours

	return fmt.Sprintf("%s🔥$%.2f/hr%s", ColorRed, rate, ColorReset)
}

// 格式化今日/週成本
func formatCostStats(daily, weekly UsageStats) string {
	dailyCostStr := formatCost(daily.TotalCost)
	weeklyCostStr := formatCost(weekly.TotalCost)
	return fmt.Sprintf("%s📆%s%s/%s💵%s%s", ColorGold, dailyCostStr, ColorReset, ColorBlue, weeklyCostStr, ColorReset)
}

// 格式化 Cache 命中率
func formatCacheHitRate(usage SessionUsageResult) string {
	totalInput := usage.InputTokens + usage.CacheReadTokens
	if totalInput == 0 {
		return fmt.Sprintf("%s📦--%s", ColorDim, ColorReset)
	}

	hitRate := float64(usage.CacheReadTokens) * 100.0 / float64(totalInput)

	// 根據命中率選擇顏色
	var color string
	if hitRate >= 70 {
		color = ColorGreen
	} else if hitRate >= 40 {
		color = ColorYellow
	} else {
		color = ColorOrange
	}

	return fmt.Sprintf("%s📦%.0f%%%s", color, hitRate, ColorReset)
}

// 分析 Context 使用量
func analyzeContext(transcriptPath string) string {
	var contextLength int

	if transcriptPath == "" {
		contextLength = 0
	} else {
		contextLength = calculateContextUsage(transcriptPath)
	}

	percentage := int(float64(contextLength) * 100.0 / 200000.0)
	if percentage > 100 {
		percentage = 100
	}

	progressBar := generateProgressBar(percentage)
	formattedNum := formatNumber(contextLength)
	color := getContextColor(percentage)

	return fmt.Sprintf(" | %s %s%d%% %s%s",
		progressBar, color, percentage, formattedNum, ColorReset)
}

// 計算 Context 使用量
func calculateContextUsage(transcriptPath string) int {
	file, err := os.Open(transcriptPath)
	if err != nil {
		return 0
	}
	defer file.Close()

	lines := make([]string, 0, 100)
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
	lines = allLines[start:]

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

// 生成進度條
func generateProgressBar(percentage int) string {
	width := 10
	filled := percentage * width / 100
	if filled > width {
		filled = width
	}

	empty := width - filled
	color := getContextColor(percentage)

	var bar strings.Builder

	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(ColorReset)
	}

	if empty > 0 {
		bar.WriteString(ColorGray)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(ColorReset)
	}

	return bar.String()
}

// 獲取 Context 顏色
func getContextColor(percentage int) string {
	if percentage < 60 {
		return ColorCtxGreen
	} else if percentage < 80 {
		return ColorCtxGold
	}
	return ColorCtxRed
}

// 格式化數字
func formatNumber(num int) string {
	if num == 0 {
		return "--"
	}

	if num >= 1000000 {
		return fmt.Sprintf("%dM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%dk", num/1000)
	}
	return strconv.Itoa(num)
}
