package themes

import (
	"fmt"
	"strings"
)

// PixelTheme 8-bit 像素遊戲風格
type PixelTheme struct{}

func init() {
	RegisterTheme(&PixelTheme{})
}

func (t *PixelTheme) Name() string {
	return "pixel"
}

func (t *PixelTheme) Description() string {
	return "像素風：8-bit 復古遊戲，方塊字符"
}

const (
	PixelRed     = "\033[38;2;255;0;77m"
	PixelOrange  = "\033[38;2;255;163;0m"
	PixelYellow  = "\033[38;2;255;236;39m"
	PixelGreen   = "\033[38;2;0;228;54m"
	PixelCyan    = "\033[38;2;41;173;255m"
	PixelBlue    = "\033[38;2;131;118;156m"
	PixelPink    = "\033[38;2;255;119;168m"
	PixelPeach   = "\033[38;2;255;204;170m"
	PixelWhite   = "\033[38;2;255;241;232m"
	PixelGray    = "\033[38;2;95;87;79m"
	PixelDark    = "\033[38;2;41;44;45m"
	PixelBgGreen = "\033[48;2;0;80;30m"
	PixelBgRed   = "\033[48;2;100;0;30m"
)

func (t *PixelTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Top border (pixel blocks)
	sb.WriteString(PixelCyan + "█" + PixelGreen + "█" + PixelYellow + "█" + PixelOrange + "█" + PixelRed + "█" + PixelPink + "█")
	sb.WriteString(PixelGray + strings.Repeat("▀", 70))
	sb.WriteString(PixelPink + "█" + PixelRed + "█" + PixelOrange + "█" + PixelYellow + "█" + PixelGreen + "█" + PixelCyan + "█")
	sb.WriteString(Reset + "\n")

	// Model + path like game HUD
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s★NEW%s", PixelYellow, Reset)
	}

	// Player info style
	line1 := fmt.Sprintf(" %s♦%s %s%s%s%s%s %s%s%s%s  %s▪%s  %s◆%s %s%s",
		PixelCyan, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		PixelGray, data.Version, Reset, update,
		PixelGray, Reset,
		PixelYellow, Reset, ShortenPath(data.ProjectPath, 20), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⬡%s%s", PixelGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", PixelGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s*%d%s", PixelOrange, data.GitDirty, Reset)
		}
	}
	sb.WriteString(line1)
	sb.WriteString("\n")

	// Game stats style (score/coins/time)
	line2 := fmt.Sprintf(" %s♦%s %sTOK%s%6s  %sMSG%s%4d  %sTIME%s%s  %s▪%s  %s%s%s  %s%s%s  %s%s%s/h  %sHIT%s%d%%",
		PixelCyan, Reset,
		PixelGray, PixelWhite, FormatTokens(data.TokenCount),
		PixelGray, PixelCyan, data.MessageCount,
		PixelGray, PixelPeach, data.SessionTime,
		PixelGray, Reset,
		PixelGreen, FormatCostShort(data.SessionCost), Reset,
		PixelYellow, FormatCostShort(data.DayCost), Reset,
		PixelRed, FormatCostShort(data.BurnRate), Reset,
		PixelGray, PixelGreen, data.CacheHitRate)
	sb.WriteString(line2)
	sb.WriteString("\n")

	// Health/Mana/XP bars (game style)
	ctxBar := t.generatePixelBar(data.ContextPercent, 12, PixelCyan, PixelBgGreen)
	bar5 := t.generatePixelBar(data.API5hrPercent, 8, PixelGreen, PixelBgGreen)
	bar7 := t.generatePixelBar(data.API7dayPercent, 8, PixelYellow, PixelBgRed)

	ctxColor := PixelGreen
	if data.ContextPercent >= 80 {
		ctxColor = PixelRed
	} else if data.ContextPercent >= 60 {
		ctxColor = PixelYellow
	}

	line3 := fmt.Sprintf(" %s♦%s %sCTX%s%s%s%3d%%%s  %s5HR%s%s%s%3d%%%s %s%s%s  %s7DY%s%s%s%3d%%%s %s%s%s",
		PixelCyan, Reset,
		PixelGray, Reset, ctxBar, ctxColor, data.ContextPercent, Reset,
		PixelGray, Reset, bar5, PixelGreen, data.API5hrPercent, Reset,
		PixelGray, data.API5hrTimeLeft, Reset,
		PixelGray, Reset, bar7, PixelYellow, data.API7dayPercent, Reset,
		PixelGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line3)
	sb.WriteString("\n")

	// Bottom border
	sb.WriteString(PixelCyan + "█" + PixelGreen + "█" + PixelYellow + "█" + PixelOrange + "█" + PixelRed + "█" + PixelPink + "█")
	sb.WriteString(PixelGray + strings.Repeat("▄", 70))
	sb.WriteString(PixelPink + "█" + PixelRed + "█" + PixelOrange + "█" + PixelYellow + "█" + PixelGreen + "█" + PixelCyan + "█")
	sb.WriteString(Reset + "\n")

	return sb.String()
}

func (t *PixelTheme) generatePixelBar(percent, width int, color, bgColor string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	bar.WriteString("〔")
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(PixelDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString("〕")
	return bar.String()
}
