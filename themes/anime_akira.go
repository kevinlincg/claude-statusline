package themes

import (
	"fmt"
	"strings"
)

// AkiraTheme Akira Neo-Tokyo style
type AkiraTheme struct{}

func init() {
	RegisterTheme(&AkiraTheme{})
}

func (t *AkiraTheme) Name() string {
	return "akira"
}

func (t *AkiraTheme) Description() string {
	return "Akira: Neo-Tokyo psychic warning interface"
}

const (
	AkiraRed    = "\033[38;2;220;20;20m"
	AkiraBlue   = "\033[38;2;0;100;200m"
	AkiraWhite  = "\033[38;2;240;240;240m"
	AkiraYellow = "\033[38;2;255;200;0m"
	AkiraCyan   = "\033[38;2;0;180;180m"
	AkiraGray   = "\033[38;2;100;100;100m"
	AkiraDark   = "\033[38;2;30;30;30m"
	AkiraPink   = "\033[38;2;255;100;150m"
)

func (t *AkiraTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(AkiraRed + "╔═══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(AkiraRed + "║" + Reset + "  " + AkiraRed + "▲ WARNING ▲" + Reset + "   " + AkiraWhite + "NEO-TOKYO ESPER MONITORING SYSTEM" + Reset + "   " + AkiraYellow + "アキラ" + Reset + "             " + AkiraRed + "║" + Reset + "\n")
	sb.WriteString(AkiraRed + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	subject := "Subject #41"
	if data.ModelType == "Opus" {
		subject = "Subject #28 AKIRA"
	} else if data.ModelType == "Haiku" {
		subject = "Subject #27"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[AWAKENING]%s", AkiraRed, Reset)
	}

	line1 := fmt.Sprintf("  %sSubject:%s %s%s%s  %sID:%s %s%s%s  %s%s%s%s",
		AkiraRed, Reset, modelColor, modelIcon, data.ModelName,
		AkiraGray, Reset, AkiraCyan, subject, Reset,
		AkiraGray, data.Version, Reset, update)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line1, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s◈%s%s", AkiraBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", AkiraCyan, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s!%d%s", AkiraRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sLocation:%s %s%s",
		AkiraBlue, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line2, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	sb.WriteString(AkiraRed + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	powerColor := AkiraBlue
	if data.ContextPercent > 75 {
		powerColor = AkiraRed
	} else if data.ContextPercent > 50 {
		powerColor = AkiraYellow
	}

	line3 := fmt.Sprintf("  %sPSYCHIC LEVEL%s  %s  %s%3d%%%s %s",
		AkiraPink, Reset, t.generateAkiraBar(data.ContextPercent, 16, powerColor),
		powerColor, data.ContextPercent, Reset,
		func() string {
			if data.ContextPercent > 80 {
				return AkiraRed + "DANGER" + Reset
			}
			return ""
		}())

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line3, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sCONTAINMENT%s    %s  %s%3d%%%s  %s%s%s",
		AkiraBlue, Reset, t.generateAkiraBar(100-data.API5hrPercent, 16, AkiraBlue),
		AkiraBlue, 100-data.API5hrPercent, Reset, AkiraGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line4, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sSUPPRESSION%s    %s  %s%3d%%%s  %s%s%s",
		AkiraYellow, Reset, t.generateAkiraBar(100-data.API7dayPercent, 16, AkiraYellow),
		AkiraYellow, 100-data.API7dayPercent, Reset, AkiraGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line5, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	sb.WriteString(AkiraRed + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sOutput:%s %s%s%s  %sTime:%s %s  %sEvents:%s %s%d%s  %sCost:%s %s%s%s",
		AkiraWhite, Reset, AkiraWhite, FormatTokens(data.TokenCount), Reset,
		AkiraGray, Reset, data.SessionTime,
		AkiraGray, Reset, AkiraCyan, data.MessageCount, Reset,
		AkiraYellow, Reset, AkiraYellow, FormatCost(data.SessionCost), Reset)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line6, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sStability:%s %s%d%%%s",
		AkiraBlue, Reset, AkiraBlue, FormatCost(data.DayCost), Reset,
		AkiraRed, Reset, AkiraRed, FormatCost(data.BurnRate), Reset,
		AkiraCyan, Reset, AkiraCyan, data.CacheHitRate, Reset)

	sb.WriteString(AkiraRed + "║" + Reset)
	sb.WriteString(PadRight(line7, 89))
	sb.WriteString(AkiraRed + "║" + Reset + "\n")

	sb.WriteString(AkiraRed + "╚═══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *AkiraTheme) generateAkiraBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(AkiraDark + "【" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▮", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(AkiraDark)
		bar.WriteString(strings.Repeat("▯", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(AkiraDark + "】" + Reset)
	return bar.String()
}
