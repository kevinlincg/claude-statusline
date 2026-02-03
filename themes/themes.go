package themes

import (
	"fmt"
	"strings"
)

// ANSI 顏色定義
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Dim    = "\033[2m"

	// 基本顏色
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

	// 亮色版本
	ColorBrightGreen  = "\033[38;2;80;255;100m"
	ColorBrightCyan   = "\033[38;2;0;255;255m"
	ColorBrightYellow = "\033[38;2;255;220;60m"
	ColorNeonGreen    = "\033[38;2;0;255;136m"
	ColorNeonPink     = "\033[38;2;255;0;255m"
	ColorNeonOrange   = "\033[38;2;255;150;50m"

	// Context 顏色
	ColorCtxGreen = "\033[38;2;108;167;108m"
	ColorCtxGold  = "\033[38;2;188;155;83m"
	ColorCtxRed   = "\033[38;2;185;102;82m"

	// 框線顏色
	ColorFrame    = "\033[38;2;60;60;60m"
	ColorFrameDim = "\033[38;2;50;50;50m"
	ColorLabel    = "\033[38;2;140;140;140m"
	ColorLabelDim = "\033[38;2;100;100;100m"
	ColorTreeDim  = "\033[38;2;100;100;100m"

	// 光棒背景色
	BgGreenGlow  = "\033[48;2;20;55;25m"
	BgYellowGlow = "\033[48;2;55;50;15m"
	BgCyanGlow   = "\033[48;2;0;60;60m"
	BgRedGlow    = "\033[48;2;60;20;20m"
)

// StatusData 包含所有要顯示的狀態資料
type StatusData struct {
	// 模型資訊
	ModelName    string
	ModelType    string // Opus, Sonnet, Haiku
	ModelIcon    string
	ModelColor   string

	// 版本資訊
	Version       string
	UpdateAvailable bool

	// 工作區資訊
	ProjectPath string
	GitBranch   string
	GitStaged   int
	GitDirty    int

	// Session 統計
	TokenCount   int64
	MessageCount int
	SessionTime  string
	CacheHitRate int

	// 成本
	SessionCost float64
	DayCost     float64
	MonthCost   float64
	WeekCost    float64
	BurnRate    float64

	// Context
	ContextUsed    int
	ContextPercent int

	// API 限制
	API5hrPercent   int
	API5hrTimeLeft  string
	API7dayPercent  int
	API7dayTimeLeft string
}

// Theme 介面定義
type Theme interface {
	Name() string
	Description() string
	Render(data StatusData) string
}

// ThemeRegistry 主題註冊表
var ThemeRegistry = make(map[string]Theme)

// RegisterTheme 註冊主題
func RegisterTheme(theme Theme) {
	ThemeRegistry[theme.Name()] = theme
}

// GetTheme 獲取主題
func GetTheme(name string) (Theme, bool) {
	theme, ok := ThemeRegistry[name]
	return theme, ok
}

// ListThemes 列出所有主題
func ListThemes() []Theme {
	themes := make([]Theme, 0, len(ThemeRegistry))
	for _, theme := range ThemeRegistry {
		themes = append(themes, theme)
	}
	return themes
}

// 輔助函數

// FormatTokens 格式化 token 數量
func FormatTokens(tokens int64) string {
	if tokens >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
	} else if tokens >= 1000 {
		return fmt.Sprintf("%.1fk", float64(tokens)/1000)
	}
	return fmt.Sprintf("%d", tokens)
}

// FormatTokensFixed 格式化 token 數量（固定寬度）
func FormatTokensFixed(tokens int64, width int) string {
	s := FormatTokens(tokens)
	return PadLeft(s, width)
}

// FormatCost 格式化成本
func FormatCost(cost float64) string {
	if cost >= 100 {
		return fmt.Sprintf("$%.0f", cost)
	} else if cost >= 10 {
		return fmt.Sprintf("$%.1f", cost)
	} else if cost >= 1 {
		return fmt.Sprintf("$%.2f", cost)
	}
	return fmt.Sprintf("$%.2f", cost)
}

// FormatCostShort 格式化成本（簡短）
func FormatCostShort(cost float64) string {
	if cost >= 100 {
		return fmt.Sprintf("$%.0f", cost)
	} else if cost >= 10 {
		return fmt.Sprintf("$%.0f", cost)
	}
	return fmt.Sprintf("$%.2f", cost)
}

// FormatPercent 格式化百分比
func FormatPercent(pct int) string {
	return fmt.Sprintf("%d%%", pct)
}

// FormatPercentFixed 格式化百分比（固定寬度）
func FormatPercentFixed(pct int, width int) string {
	s := fmt.Sprintf("%d%%", pct)
	return PadLeft(s, width)
}

// FormatNumber 格式化數字
func FormatNumber(num int) string {
	if num >= 1000000 {
		return fmt.Sprintf("%dM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%dk", num/1000)
	}
	return fmt.Sprintf("%d", num)
}

// GenerateBar 生成進度條
func GenerateBar(percent, width int, filledChar, emptyChar string, filledColor, emptyColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(filledColor)
		bar.WriteString(strings.Repeat(filledChar, filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(emptyColor)
		bar.WriteString(strings.Repeat(emptyChar, empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}

// GenerateGlowBar 生成發光進度條
func GenerateGlowBar(percent, width int, color, bgColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(bgColor)
		bar.WriteString(Bold)
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString("\033[38;2;35;35;35m")
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}

// GetBarColor 根據百分比獲取顏色
func GetBarColor(percent int) (string, string) {
	if percent < 50 {
		return ColorBrightGreen, BgGreenGlow
	} else if percent < 75 {
		return ColorBrightYellow, BgYellowGlow
	}
	return ColorRed, BgRedGlow
}

// GetContextColor 根據 context 百分比獲取顏色
func GetContextColor(percent int) string {
	if percent < 60 {
		return ColorCtxGreen
	} else if percent < 80 {
		return ColorCtxGold
	}
	return ColorCtxRed
}

// PadLeft 左填充
func PadLeft(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	return strings.Repeat(" ", width-visible) + s
}

// PadRight 右填充
func PadRight(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visible)
}

// PadCenter 置中填充
func PadCenter(s string, width int) string {
	visible := VisibleWidth(s)
	if visible >= width {
		return s
	}
	left := (width - visible) / 2
	right := width - visible - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// VisibleWidth 計算可見寬度（排除 ANSI 碼）
func VisibleWidth(s string) int {
	// 移除 ANSI escape codes
	clean := s
	for {
		start := strings.Index(clean, "\033[")
		if start == -1 {
			break
		}
		end := strings.Index(clean[start:], "m")
		if end == -1 {
			break
		}
		clean = clean[:start] + clean[start+end+1:]
	}

	width := 0
	for _, r := range clean {
		w := RuneWidth(r)
		width += w
	}
	return width
}

// RuneWidth 計算單個 rune 的顯示寬度
func RuneWidth(r rune) int {
	// Variation selectors - zero width
	if r >= 0xFE00 && r <= 0xFE0F {
		return 0
	}
	// Zero-width characters
	if r == 0x200B || r == 0x200C || r == 0x200D || r == 0xFEFF {
		return 0
	}
	// Combining characters - zero width
	if r >= 0x0300 && r <= 0x036F {
		return 0
	}

	// Wide characters (2 cells)
	if r >= 0x1F300 && r <= 0x1FAFF {
		return 2
	}
	if r >= 0x2300 && r <= 0x23FF {
		return 2
	}
	if r >= 0x2600 && r <= 0x26FF {
		return 2
	}
	if r >= 0x2700 && r <= 0x27BF {
		return 2
	}
	if r >= 0x2B50 && r <= 0x2B55 {
		return 2
	}
	if r >= 0x4E00 && r <= 0x9FFF {
		return 2
	}
	if r >= 0x3000 && r <= 0x303F {
		return 2
	}
	if r >= 0xFF00 && r <= 0xFFEF {
		return 2
	}

	return 1
}

// GetModelConfig 獲取模型配置
func GetModelConfig(modelType string) (color string, icon string) {
	switch modelType {
	case "Opus":
		return ColorGold, "💛"
	case "Sonnet":
		return ColorCyan, "💠"
	case "Haiku":
		return ColorPink, "🌸"
	default:
		return ColorCyan, "◆"
	}
}
