package themes

import (
	"fmt"
	"strings"
)

// GlitchTheme 數位故障風格
type GlitchTheme struct{}

func init() {
	RegisterTheme(&GlitchTheme{})
}

func (t *GlitchTheme) Name() string {
	return "glitch"
}

func (t *GlitchTheme) Description() string {
	return "故障風：數位錯位，賽博龐克破碎美學"
}

const (
	GlitchRed   = "\033[38;2;255;0;60m"
	GlitchCyan  = "\033[38;2;0;255;240m"
	GlitchWhite = "\033[38;2;255;255;255m"
	GlitchGray  = "\033[38;2;80;80;80m"
	GlitchDim   = "\033[38;2;50;50;50m"
	GlitchPink  = "\033[38;2;255;60;150m"
)

func (t *GlitchTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Glitch top border
	border := GlitchDim + "▓▒░" + GlitchRed + "█" + GlitchDim + "░" + GlitchGray + strings.Repeat("▀", 68) + GlitchDim + "░" + GlitchCyan + "█" + GlitchDim + "░▒▓" + Reset
	sb.WriteString(border + "\n")

	// Model info line
	modelColor, _ := GetModelConfig(data.ModelType)
	modelStr := fmt.Sprintf("%s%s%s", modelColor, data.ModelName, Reset)
	verStr := fmt.Sprintf("%s%s%s", GlitchGray, data.Version, Reset)

	update := ""
	if data.UpdateAvailable {
		update = GlitchRed + " [!]" + Reset
	}

	pathStr := fmt.Sprintf("%s%s%s", GlitchWhite, ShortenPath(data.ProjectPath, 25), Reset)
	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf(" %s<%s>%s", GlitchCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf(" %s+%d%s", GlitchCyan, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf(" %s~%d%s", GlitchRed, data.GitDirty, Reset)
		}
	}

	line1 := fmt.Sprintf("%s▌%s %s %s%s  %s▐▌%s  %s%s",
		GlitchRed, Reset,
		modelStr, verStr, update,
		GlitchCyan, Reset,
		pathStr, gitStr)
	sb.WriteString(line1 + "\n")

	// Stats line with chromatic split effect
	line2 := fmt.Sprintf("%s▌%s %sTOK%s %s%-6s%s  %sMSG%s %s%-3d%s  %sTIME%s %s%-6s%s  %s▐▌%s  %sSES%s %s%s%s  %sDAY%s %s%s%s  %sRATE%s %s%s/h%s",
		GlitchRed, Reset,
		GlitchGray, Reset, GlitchPink, FormatTokens(data.TokenCount), Reset,
		GlitchGray, Reset, GlitchCyan, data.MessageCount, Reset,
		GlitchGray, Reset, GlitchWhite, data.SessionTime, Reset,
		GlitchCyan, Reset,
		GlitchGray, Reset, GlitchCyan, FormatCostShort(data.SessionCost), Reset,
		GlitchGray, Reset, GlitchWhite, FormatCostShort(data.DayCost), Reset,
		GlitchGray, Reset, GlitchRed, FormatCostShort(data.BurnRate), Reset)
	sb.WriteString(line2 + "\n")

	// Progress bars
	ctxBar := t.generateGlitchBar(data.ContextPercent, 12)
	bar5 := t.generateGlitchBar(data.API5hrPercent, 10)
	bar7 := t.generateGlitchBar(data.API7dayPercent, 10)

	ctxColor := GlitchCyan
	if data.ContextPercent >= 80 {
		ctxColor = GlitchRed
	} else if data.ContextPercent >= 60 {
		ctxColor = GlitchPink
	}

	line3 := fmt.Sprintf("%s▌%s %sCTX%s%s%s%3d%%%s  %s5HR%s%s%s%3d%%%s %s%-5s%s  %s7DY%s%s%s%3d%%%s %s%-5s%s  %sHIT%s %s%d%%%s",
		GlitchRed, Reset,
		GlitchGray, Reset, ctxBar, ctxColor, data.ContextPercent, Reset,
		GlitchGray, Reset, bar5, GlitchCyan, data.API5hrPercent, Reset,
		GlitchGray, data.API5hrTimeLeft, Reset,
		GlitchGray, Reset, bar7, GlitchPink, data.API7dayPercent, Reset,
		GlitchGray, data.API7dayTimeLeft, Reset,
		GlitchGray, Reset, GlitchCyan, data.CacheHitRate, Reset)
	sb.WriteString(line3 + "\n")

	// Bottom border
	borderBot := GlitchDim + "▓▒░" + GlitchCyan + "█" + GlitchDim + "░" + GlitchGray + strings.Repeat("▄", 68) + GlitchDim + "░" + GlitchRed + "█" + GlitchDim + "░▒▓" + Reset
	sb.WriteString(borderBot + "\n")

	return sb.String()
}

func (t *GlitchTheme) generateGlitchBar(percent, width int) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(GlitchDim + "[" + Reset)

	if filled > 0 {
		for i := 0; i < filled; i++ {
			if i == filled/2 && filled > 3 {
				bar.WriteString(GlitchRed + "#" + Reset)
			} else {
				bar.WriteString(GlitchCyan + "=" + Reset)
			}
		}
	}
	if empty > 0 {
		bar.WriteString(GlitchDim)
		bar.WriteString(strings.Repeat("-", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(GlitchDim + "]" + Reset)
	return bar.String()
}
