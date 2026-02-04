package themes

import (
	"fmt"
	"strings"
)

// CyberpunkTheme 賽博朋克霓虹風格
type CyberpunkTheme struct{}

func init() {
	RegisterTheme(&CyberpunkTheme{})
}

func (t *CyberpunkTheme) Name() string {
	return "cyberpunk"
}

func (t *CyberpunkTheme) Description() string {
	return "賽博朋克：霓虹雙色框線"
}

const (
	CyberCyan    = "\033[38;2;0;255;255m"
	CyberMagenta = "\033[38;2;255;0;255m"
)

func (t *CyberpunkTheme) Render(data StatusData) string {
	var sb strings.Builder

	const width = 95

	// 頂部框線
	sb.WriteString(CyberCyan)
	sb.WriteString("╔")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╗")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第一行：模型 + 版本 | 路徑 + Git
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorNeonOrange, Reset)
	}

	line1 := fmt.Sprintf(" %s%s%s%s%s %s%s%s%s  %s│%s  %s📂 %s%s",
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		ColorNeonGreen, data.Version, Reset, update,
		ColorDim, Reset,
		ColorYellow, data.ProjectPath, Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⚡%s%s", CyberCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
	}

	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line1, width))
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 分隔線
	sb.WriteString(CyberMagenta)
	sb.WriteString("╠")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╣")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第二行：Session 統計 | Cost
	line2 := fmt.Sprintf(" %s%5s%s tok  %s%3d%s msg  %s%6s%s  %s│%s  %s%s%s ses  %s%s%s day  %s%s%s mon  %s%s/h%s  %s%d%%hit%s",
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset,
		ColorDim, Reset,
		ColorGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorPurple, FormatCost(data.MonthCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset)

	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line2, width))
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 第三行：光棒
	color1, bg1 := GetBarColor(data.ContextPercent)
	color5, bg5 := GetBarColor(data.API5hrPercent)
	color7, bg7 := GetBarColor(data.API7dayPercent)

	line3 := fmt.Sprintf(" %sCtx%s %s %s%3d%%%s  %s│%s  %s5hr%s %s %s%3d%%%s %s%s%s  %s│%s  %s7dy%s %s %s%3d%%%s %s%s%s",
		ColorLabelDim, Reset,
		GenerateGlowBar(data.ContextPercent, 15, color1, bg1),
		color1, data.ContextPercent, Reset,
		ColorDim, Reset,
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API5hrPercent, 10, color5, bg5),
		color5, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset,
		ColorDim, Reset,
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API7dayPercent, 10, color7, bg7),
		color7, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)

	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line3, width))
	sb.WriteString(CyberCyan)
	sb.WriteString("║")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// 底部框線
	sb.WriteString(CyberMagenta)
	sb.WriteString("╚")
	sb.WriteString(strings.Repeat("═", width))
	sb.WriteString("╝")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}
