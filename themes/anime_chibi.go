package themes

import (
	"fmt"
	"strings"
)

// ChibiTheme Chibi cute compact style
type ChibiTheme struct{}

func init() {
	RegisterTheme(&ChibiTheme{})
}

func (t *ChibiTheme) Name() string {
	return "chibi"
}

func (t *ChibiTheme) Description() string {
	return "Chibi: Kawaii super-deformed compact style"
}

const (
	CHBPink   = "\033[38;2;255;150;200m"
	CHBBlue   = "\033[38;2;150;200;255m"
	CHBYellow = "\033[38;2;255;230;100m"
	CHBGreen  = "\033[38;2;150;230;150m"
	CHBWhite  = "\033[38;2;255;255;255m"
	CHBGray   = "\033[38;2;180;180;180m"
)

func (t *ChibiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Cute bouncy header
	sb.WriteString("\n")
	sb.WriteString("    " + CHBPink + "(◕‿◕)" + Reset + " " + CHBYellow + "★" + CHBPink + "･" + CHBBlue + "｡" + CHBGreen + "ﾟ" + CHBYellow + "★" + Reset + " ")
	sb.WriteString(CHBWhite + "C H I B I   M O D E" + Reset + " ")
	sb.WriteString(CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + Reset + " " + CHBBlue + "(◕‿◕)" + Reset + "\n")
	sb.WriteString("                            " + CHBPink + "ちび" + Reset + "\n")
	sb.WriteString("\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	chibi := "(｡◕‿◕｡)"
	if data.ModelType == "Opus" {
		chibi = "(◕ᴗ◕✿)"
	} else if data.ModelType == "Haiku" {
		chibi = "(◠‿◠)"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s★NEW★%s", CHBYellow, Reset)
	}

	line1 := fmt.Sprintf("  %s%s%s %s%s%s %s%s%s%s",
		CHBPink, chibi, Reset, modelColor, modelIcon, data.ModelName,
		CHBGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf(" %s♪%s%s", CHBBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", CHBGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", CHBPink, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %s♡%s %s%s",
		CHBPink, Reset, ShortenPath(data.ProjectPath, 50), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString("\n")

	// Compact cute bars
	contextColor := CHBPink
	if data.ContextPercent > 75 {
		contextColor = CHBYellow
	}

	line3 := fmt.Sprintf("  %s(ﾉ◕ヮ◕)ﾉ%s %s %s%2d%%%s",
		contextColor, Reset, t.generateCHBBar(data.ContextPercent, 12, contextColor), contextColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %s٩(◕‿◕｡)۶%s %s %s%2d%%%s %s%s%s",
		CHBBlue, Reset, t.generateCHBBar(100-data.API5hrPercent, 12, CHBBlue),
		CHBBlue, 100-data.API5hrPercent, Reset, CHBGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %s(◕‿◕)♡%s  %s %s%2d%%%s %s%s%s",
		CHBGreen, Reset, t.generateCHBBar(100-data.API7dayPercent, 12, CHBGreen),
		CHBGreen, 100-data.API7dayPercent, Reset, CHBGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString("\n")

	line6 := fmt.Sprintf("  %s%s%s %s%s%s %s%d%s %s$%s%s %s%d%%%s",
		CHBWhite, FormatTokens(data.TokenCount), Reset,
		CHBGray, data.SessionTime, Reset,
		CHBBlue, data.MessageCount, Reset,
		CHBYellow, FormatCost(data.SessionCost), Reset,
		CHBPink, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString("\n")
	sb.WriteString("    " + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + CHBBlue + "｡" + CHBPink + "･" + CHBYellow + "★" + CHBGreen + "ﾟ" + Reset + "\n")

	return sb.String()
}

func (t *ChibiTheme) generateCHBBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(CHBGray + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("♥", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(CHBGray)
		bar.WriteString(strings.Repeat("♡", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(CHBGray + "〉" + Reset)
	return bar.String()
}
