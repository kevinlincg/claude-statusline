package themes

import (
	"fmt"
	"strings"
)

// HtopTheme htop 系統監視器風格
type HtopTheme struct{}

func init() {
	RegisterTheme(&HtopTheme{})
}

func (t *HtopTheme) Name() string {
	return "htop"
}

func (t *HtopTheme) Description() string {
	return "htop：經典系統監視器，彩色進度條風格"
}

const (
	HtopBlack   = "\033[38;2;0;0;0m"
	HtopRed     = "\033[38;2;205;0;0m"
	HtopGreen   = "\033[38;2;0;205;0m"
	HtopYellow  = "\033[38;2;205;205;0m"
	HtopBlue    = "\033[38;2;0;0;238m"
	HtopMagenta = "\033[38;2;205;0;205m"
	HtopCyan    = "\033[38;2;0;205;205m"
	HtopWhite   = "\033[38;2;229;229;229m"
	HtopBrightBlack  = "\033[38;2;127;127;127m"
	HtopBrightRed    = "\033[38;2;255;0;0m"
	HtopBrightGreen  = "\033[38;2;0;255;0m"
	HtopBrightYellow = "\033[38;2;255;255;0m"
	HtopBrightBlue   = "\033[38;2;92;92;255m"
	HtopBrightCyan   = "\033[38;2;0;255;255m"
	HtopBrightWhite  = "\033[38;2;255;255;255m"
	HtopBgBlue       = "\033[48;2;0;0;138m"
	HtopBgCyan       = "\033[48;2;0;139;139m"
	HtopBgBlack      = "\033[48;2;0;0;0m"
)

func (t *HtopTheme) Render(data StatusData) string {
	var sb strings.Builder

	// CPU-style meters header
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	// CPU bars (context as CPU usage style)
	cpuUsed := data.ContextPercent
	cpu1 := t.generateHtopMeter(cpuUsed, 25, "0")
	cpu2 := t.generateHtopMeter(data.API5hrPercent, 25, "1")

	line1 := fmt.Sprintf("%s0%s[%s%s%s]%s  %s1%s[%s%s%s]%s  %s%s%s%s %s%s%s",
		HtopBrightCyan, HtopBrightBlack, Reset, cpu1, HtopBrightBlack, Reset,
		HtopBrightCyan, HtopBrightBlack, Reset, cpu2, HtopBrightBlack, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		HtopBrightBlack, data.Version)
	if data.UpdateAvailable {
		line1 += HtopBrightYellow + " [UPDATE]" + Reset
	}
	sb.WriteString(line1 + "\n")

	// Memory-style bar
	memBar := t.generateHtopMemBar(data.API7dayPercent, 25)
	line2 := fmt.Sprintf("%sMem%s[%s%s%s]%s  %sTasks:%s %s%d%s  %sLoad:%s %s%s%s  %sUptime:%s %s%s%s",
		HtopBrightGreen, HtopBrightBlack, Reset, memBar, HtopBrightBlack, Reset,
		HtopBrightBlack, Reset, HtopBrightWhite, data.MessageCount, Reset,
		HtopBrightBlack, Reset, HtopBrightWhite, FormatTokens(data.TokenCount), Reset,
		HtopBrightBlack, Reset, HtopBrightWhite, data.SessionTime, Reset)
	sb.WriteString(line2 + "\n")

	// Swap-style bar (burn rate indicator)
	swpBar := t.generateHtopSwapBar(data.CacheHitRate, 25)
	line3 := fmt.Sprintf("%sSwp%s[%s%s%s]%s  %sPath:%s %s%s%s",
		HtopBrightRed, HtopBrightBlack, Reset, swpBar, HtopBrightBlack, Reset,
		HtopBrightBlack, Reset, HtopCyan, ShortenPath(data.ProjectPath, 35), Reset)
	if data.GitBranch != "" {
		line3 += fmt.Sprintf("  %s<%s>%s", HtopBrightGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line3 += fmt.Sprintf(" %s+%d%s", HtopGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line3 += fmt.Sprintf(" %s~%d%s", HtopYellow, data.GitDirty, Reset)
		}
	}
	sb.WriteString(line3 + "\n")

	// Separator
	sb.WriteString(HtopBrightBlack + strings.Repeat("─", 80) + Reset + "\n")

	// Process-style info line
	line4 := fmt.Sprintf("  %sPID%s  %sUSER%s      %sCPU%%%s  %sMEM%%%s  %sTIME+%s     %sCOMMAND%s",
		HtopBrightBlack, Reset, HtopBrightBlack, Reset,
		HtopBrightBlack, Reset, HtopBrightBlack, Reset,
		HtopBrightBlack, Reset, HtopBrightBlack, Reset)
	sb.WriteString(line4 + "\n")

	// Main "process" line
	ctxPct := 100 - data.ContextPercent
	apiPct := 100 - data.API5hrPercent
	line5 := fmt.Sprintf("  %s%-5s%s %s%-9s%s %s%5.1f%s  %s%5.1f%s  %s%-10s%s %sclaude-session%s",
		HtopBrightCyan, "1", Reset,
		HtopBrightGreen, "claude", Reset,
		HtopBrightRed, float64(100-ctxPct), Reset,
		HtopBrightGreen, float64(100-apiPct), Reset,
		HtopBrightWhite, data.SessionTime, Reset,
		HtopBrightWhite, Reset)
	sb.WriteString(line5 + "\n")

	// F-key menu bar (htop signature)
	fkeys := fmt.Sprintf("%sF1%sHelp %sF2%sSetup %sF3%sSearch %sF4%sFilter %sF5%sTree %sF6%sSortBy %sF7%sNice- %sF8%sNice+ %sF9%sKill %sF10%sQuit",
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack,
		HtopBgBlack+HtopBrightCyan, HtopBgCyan+HtopBlack)
	sb.WriteString(fkeys + Reset + "\n")

	return sb.String()
}

func (t *HtopTheme) generateHtopMeter(percent, width int, label string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	// htop uses gradient: green -> yellow -> red
	greenPart := filled
	yellowPart := 0
	redPart := 0
	if percent > 50 {
		greenPart = width / 2
		yellowPart = filled - greenPart
		if percent > 75 {
			yellowPart = width / 4
			redPart = filled - greenPart - yellowPart
		}
	}

	if greenPart > 0 {
		bar.WriteString(HtopBrightGreen)
		bar.WriteString(strings.Repeat("|", greenPart))
	}
	if yellowPart > 0 {
		bar.WriteString(HtopBrightYellow)
		bar.WriteString(strings.Repeat("|", yellowPart))
	}
	if redPart > 0 {
		bar.WriteString(HtopBrightRed)
		bar.WriteString(strings.Repeat("|", redPart))
	}
	if empty > 0 {
		bar.WriteString(HtopBrightBlack)
		bar.WriteString(strings.Repeat(" ", empty))
	}
	bar.WriteString(Reset)
	bar.WriteString(fmt.Sprintf("%s%3d%%%s", HtopBrightWhite, percent, Reset))
	return bar.String()
}

func (t *HtopTheme) generateHtopMemBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(HtopBrightGreen)
		bar.WriteString(strings.Repeat("|", filled))
	}
	if empty > 0 {
		bar.WriteString(HtopBrightBlack)
		bar.WriteString(strings.Repeat(" ", empty))
	}
	bar.WriteString(Reset)
	bar.WriteString(fmt.Sprintf("%s%3d%%%s", HtopBrightWhite, percent, Reset))
	return bar.String()
}

func (t *HtopTheme) generateHtopSwapBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(HtopBrightRed)
		bar.WriteString(strings.Repeat("|", filled))
	}
	if empty > 0 {
		bar.WriteString(HtopBrightBlack)
		bar.WriteString(strings.Repeat(" ", empty))
	}
	bar.WriteString(Reset)
	bar.WriteString(fmt.Sprintf("%s%3d%%%s", HtopBrightWhite, percent, Reset))
	return bar.String()
}
