package themes

import (
	"fmt"
	"strings"
)

// IdolTheme Idol anime stage concert style
type IdolTheme struct{}

func init() {
	RegisterTheme(&IdolTheme{})
}

func (t *IdolTheme) Name() string {
	return "idol"
}

func (t *IdolTheme) Description() string {
	return "Idol: Concert stage with light sticks and stars"
}

const (
	IDLPink   = "\033[38;2;255;100;150m"
	IDLBlue   = "\033[38;2;100;150;255m"
	IDLYellow = "\033[38;2;255;220;0m"
	IDLPurple = "\033[38;2;200;100;255m"
	IDLCyan   = "\033[38;2;0;220;220m"
	IDLWhite  = "\033[38;2;255;255;255m"
	IDLGray   = "\033[38;2;150;150;150m"
)

func (t *IdolTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Stage lights
	sb.WriteString("  " + IDLYellow + "✦" + IDLPink + "✧" + IDLBlue + "✦" + IDLPurple + "✧" + IDLCyan + "✦" + IDLYellow + "✧" + IDLPink + "✦" + IDLBlue + "✧" + IDLPurple + "✦" + IDLCyan + "✧" + IDLYellow + "✦" + IDLPink + "✧" + IDLBlue + "✦" + IDLPurple + "✧" + IDLCyan + "✦" + IDLYellow + "✧" + IDLPink + "✦" + IDLBlue + "✧" + IDLPurple + "✦" + IDLCyan + "✧" + IDLYellow + "✦" + IDLPink + "✧" + IDLBlue + "✦" + IDLPurple + "✧" + IDLCyan + "✦" + Reset + "\n")
	sb.WriteString("\n")

	// Stage name
	sb.WriteString("          " + IDLPink + "╱" + IDLWhite + "╲" + IDLPink + "╱" + IDLWhite + "╲" + Reset + " ")
	sb.WriteString(IDLYellow + "★" + IDLWhite + " I D O L   S T A G E " + IDLYellow + "★" + Reset + " ")
	sb.WriteString(IDLBlue + "╱" + IDLWhite + "╲" + IDLBlue + "╱" + IDLWhite + "╲" + Reset + "\n")
	sb.WriteString("                            " + IDLPink + "アイドル" + Reset + "\n")
	sb.WriteString("\n")

	// Light sticks wave
	sb.WriteString("  " + IDLPink + "│" + IDLBlue + "│" + IDLYellow + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLPink + "│" + IDLBlue + "│" + IDLYellow + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLPink + "│" + IDLBlue + "│" + IDLYellow + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLPink + "│" + IDLBlue + "│" + IDLYellow + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLPink + "│" + IDLBlue + "│" + IDLYellow + "│" + Reset + "\n")
	sb.WriteString("\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	idol := "Trainee"
	if data.ModelType == "Opus" {
		idol = "Center"
	} else if data.ModelType == "Haiku" {
		idol = "Backup"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s★DEBUT★%s", IDLYellow, Reset)
	}

	line1 := fmt.Sprintf("      %s♪ Idol:%s %s%s%s  %s♪ Position:%s %s%s%s  %s%s%s%s",
		IDLPink, Reset, modelColor, modelIcon, data.ModelName,
		IDLBlue, Reset, IDLYellow, idol, Reset,
		IDLGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s♫%s%s", IDLPurple, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", IDLPink, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", IDLPurple, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("      %s♪ Venue:%s %s%s",
		IDLCyan, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString("\n")

	// Audience cheers divider
	sb.WriteString("  " + IDLCyan + "│" + IDLPurple + "│" + IDLYellow + "│" + IDLBlue + "│" + IDLPink + "│" + IDLCyan + "│" + IDLPurple + "│" + IDLYellow + "│" + IDLBlue + "│" + IDLPink + "│" + IDLCyan + "│" + IDLPurple + "│" + IDLYellow + "│" + IDLBlue + "│" + IDLPink + "│" + IDLCyan + "│" + IDLPurple + "│" + IDLYellow + "│" + IDLBlue + "│" + IDLPink + "│" + IDLCyan + "│" + IDLPurple + "│" + IDLYellow + "│" + Reset + "\n")
	sb.WriteString("\n")

	// Stats as stage metrics
	popularityColor := IDLPink
	if data.ContextPercent > 75 {
		popularityColor = IDLPurple
	}

	line3 := fmt.Sprintf("        %s★ Popularity%s  %s  %s%3d%%%s",
		IDLPink, Reset, t.generateIDLBar(data.ContextPercent, 14, popularityColor), popularityColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("        %s★ Stamina%s     %s  %s%3d%%%s  %s%s%s",
		IDLBlue, Reset, t.generateIDLBar(100-data.API5hrPercent, 14, IDLBlue),
		IDLBlue, 100-data.API5hrPercent, Reset, IDLGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("        %s★ Fans%s        %s  %s%3d%%%s  %s%s%s",
		IDLYellow, Reset, t.generateIDLBar(100-data.API7dayPercent, 14, IDLYellow),
		IDLYellow, 100-data.API7dayPercent, Reset, IDLGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString("\n")
	sb.WriteString("  " + IDLYellow + "│" + IDLPink + "│" + IDLBlue + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLYellow + "│" + IDLPink + "│" + IDLBlue + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLYellow + "│" + IDLPink + "│" + IDLBlue + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLYellow + "│" + IDLPink + "│" + IDLBlue + "│" + IDLPurple + "│" + IDLCyan + "│" + IDLYellow + "│" + IDLPink + "│" + IDLBlue + "│" + Reset + "\n")
	sb.WriteString("\n")

	line6 := fmt.Sprintf("      %s%s%s notes  %s%s%s  %s%d%s songs  %s$%s%s  %s$%s/hr%s  %s%d%%%s",
		IDLWhite, FormatTokens(data.TokenCount), Reset,
		IDLGray, data.SessionTime, Reset,
		IDLPink, data.MessageCount, Reset,
		IDLYellow, FormatCost(data.SessionCost), Reset,
		IDLCyan, FormatCost(data.BurnRate), Reset,
		IDLPurple, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString("\n")
	sb.WriteString("  " + IDLCyan + "✧" + IDLPurple + "✦" + IDLYellow + "✧" + IDLPink + "✦" + IDLBlue + "✧" + IDLCyan + "✦" + IDLPurple + "✧" + IDLYellow + "✦" + IDLPink + "✧" + IDLBlue + "✦" + IDLCyan + "✧" + IDLPurple + "✦" + IDLYellow + "✧" + IDLPink + "✦" + IDLBlue + "✧" + IDLCyan + "✦" + IDLPurple + "✧" + IDLYellow + "✦" + IDLPink + "✧" + IDLBlue + "✦" + IDLCyan + "✧" + IDLPurple + "✦" + IDLYellow + "✧" + Reset + "\n")

	return sb.String()
}

func (t *IdolTheme) generateIDLBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(IDLGray + "〚" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("★", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(IDLGray)
		bar.WriteString(strings.Repeat("☆", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(IDLGray + "〛" + Reset)
	return bar.String()
}
