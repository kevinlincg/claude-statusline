package themes

import (
	"fmt"
	"strings"
)

// SailorMoonTheme Sailor Moon magical girl style
type SailorMoonTheme struct{}

func init() {
	RegisterTheme(&SailorMoonTheme{})
}

func (t *SailorMoonTheme) Name() string {
	return "sailormoon"
}

func (t *SailorMoonTheme) Description() string {
	return "Sailor Moon: Magical girl transformation style"
}

const (
	SMPink   = "\033[38;2;255;182;193m"
	SMYellow = "\033[38;2;255;255;150m"
	SMBlue   = "\033[38;2;135;206;250m"
	SMPurple = "\033[38;2;221;160;221m"
	SMGold   = "\033[38;2;255;215;0m"
	SMWhite  = "\033[38;2;255;255;255m"
	SMRed    = "\033[38;2;255;105;180m"
	SMGray   = "\033[38;2;180;180;180m"
)

func (t *SailorMoonTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(SMPink + "✧･ﾟ: *✧･ﾟ:*" + SMYellow + " ☾ " + SMWhite + "MOON PRISM POWER" + SMYellow + " ☾ " + SMPink + "*:･ﾟ✧*:･ﾟ✧" + Reset + "\n")
	sb.WriteString(SMPurple + "────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	guardian := "Moon"
	if data.ModelType == "Opus" {
		guardian = "Cosmos"
	} else if data.ModelType == "Haiku" {
		guardian = "Chibi Moon"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✨Transform!✨%s", SMGold, Reset)
	}

	line1 := fmt.Sprintf("  %s☆%s %sSailor%s %s%s%s  %s%s%s  %s%s%s%s",
		SMYellow, Reset,
		SMPink, Reset, SMPink, guardian, Reset,
		modelColor, modelIcon, data.ModelName,
		SMGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s♡%s%s", SMPink, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SMYellow, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SMRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %s♪%s %sMission:%s %s%s",
		SMBlue, Reset, SMBlue, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(SMPurple + "────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	moonColor := SMYellow
	if data.ContextPercent > 75 {
		moonColor = SMRed
	}

	line3 := fmt.Sprintf("  %s☾ Moon Power%s   %s  %s%3d%%%s",
		SMYellow, Reset, t.generateSMBar(data.ContextPercent, 18, moonColor), moonColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %s♡ Love Energy%s  %s  %s%3d%%%s  %s%s%s",
		SMPink, Reset, t.generateSMBar(100-data.API5hrPercent, 18, SMPink),
		SMPink, 100-data.API5hrPercent, Reset, SMGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %s★ Star Light%s   %s  %s%3d%%%s  %s%s%s",
		SMGold, Reset, t.generateSMBar(100-data.API7dayPercent, 18, SMGold),
		SMGold, 100-data.API7dayPercent, Reset, SMGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(SMPurple + "────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("  %s✧%s %s%s%s magic  %s%s%s  %s%d%s acts  %s$%s%s  %s$%s/day%s",
		SMYellow, Reset, SMPurple, FormatTokens(data.TokenCount), Reset,
		SMGray, data.SessionTime, Reset,
		SMBlue, data.MessageCount, Reset,
		SMGold, FormatCost(data.SessionCost), Reset,
		SMPink, FormatCost(data.DayCost), Reset)
	sb.WriteString(line6 + "\n")

	line7 := fmt.Sprintf("  %s♪%s Rate: %s$%s/h%s  Accuracy: %s%d%%%s %s♡%s",
		SMBlue, Reset, SMRed, FormatCost(data.BurnRate), Reset,
		SMYellow, data.CacheHitRate, Reset, SMPink, Reset)
	sb.WriteString(line7 + "\n")

	sb.WriteString(SMPink + "✧･ﾟ: *✧･ﾟ:*" + SMYellow + " ☆ " + SMWhite + "In the name of the Moon!" + SMYellow + " ☆ " + SMPink + "*:･ﾟ✧*:･ﾟ✧" + Reset + "\n")

	return sb.String()
}

func (t *SailorMoonTheme) generateSMBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SMGray + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("♥", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SMGray)
		bar.WriteString(strings.Repeat("♡", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SMGray + "〕" + Reset)
	return bar.String()
}
