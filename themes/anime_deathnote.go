package themes

import (
	"fmt"
	"strings"
)

// DeathNoteTheme Death Note notebook style
type DeathNoteTheme struct{}

func init() {
	RegisterTheme(&DeathNoteTheme{})
}

func (t *DeathNoteTheme) Name() string {
	return "deathnote"
}

func (t *DeathNoteTheme) Description() string {
	return "Death Note: Shinigami notebook gothic style"
}

const (
	DNBlack   = "\033[38;2;20;20;20m"
	DNWhite   = "\033[38;2;240;240;240m"
	DNRed     = "\033[38;2;139;0;0m"
	DNGray    = "\033[38;2;80;80;80m"
	DNDark    = "\033[38;2;40;40;40m"
	DNPurple  = "\033[38;2;75;0;130m"
	DNGold    = "\033[38;2;184;134;11m"
	DNBgBlack = "\033[48;2;10;10;10m"
)

func (t *DeathNoteTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Gothic notebook style
	sb.WriteString(DNGray + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(DNGray + "║" + Reset + "                        " + DNWhite + "D E A T H   N O T E" + Reset + "                                    " + DNGray + "║" + Reset + "\n")
	sb.WriteString(DNGray + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[SHINIGAMI EYES]%s", DNRed, Reset)
	}

	line1 := fmt.Sprintf("  %sOwner:%s %s%s%s  %s%s%s%s",
		DNPurple, Reset, modelColor, modelIcon, data.ModelName,
		DNGray, data.Version, Reset, update)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s†%s%s", DNPurple, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", DNWhite, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", DNRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sTarget:%s %s%s",
		DNRed, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	sb.WriteString(DNGray + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Rules style
	sb.WriteString(DNGray + "║" + Reset + "  " + DNWhite + "RULE I:" + Reset + "   " + DNGray + "The human whose name is written shall use context." + Reset + "              " + DNGray + "║" + Reset + "\n")

	lifeColor := DNWhite
	if data.ContextPercent > 75 {
		lifeColor = DNRed
	}

	line3 := fmt.Sprintf("  %sLifespan%s    %s  %s%3d%%%s",
		DNRed, Reset, t.generateDNBar(data.ContextPercent, 18, lifeColor), lifeColor, data.ContextPercent, Reset)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sPages%s       %s  %s%3d%%%s  %sRegen: %s%s",
		DNPurple, Reset, t.generateDNBar(100-data.API5hrPercent, 18, DNPurple),
		DNPurple, 100-data.API5hrPercent, Reset, DNDark, data.API5hrTimeLeft, Reset)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sInk%s         %s  %s%3d%%%s  %sRefill: %s%s",
		DNGold, Reset, t.generateDNBar(100-data.API7dayPercent, 18, DNGold),
		DNGold, 100-data.API7dayPercent, Reset, DNDark, data.API7dayTimeLeft, Reset)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	sb.WriteString(DNGray + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sNames:%s %s%s%s  %sTime:%s %s  %sEntries:%s %s%d%s  %sApples:%s %s%s%s",
		DNWhite, Reset, DNWhite, FormatTokens(data.TokenCount), Reset,
		DNGray, Reset, data.SessionTime,
		DNGray, Reset, DNPurple, data.MessageCount, Reset,
		DNRed, Reset, DNRed, FormatCost(data.SessionCost), Reset)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sAccuracy:%s %s%d%%%s",
		DNPurple, Reset, DNPurple, FormatCost(data.DayCost), Reset,
		DNRed, Reset, DNRed, FormatCost(data.BurnRate), Reset,
		DNGold, Reset, DNGold, data.CacheHitRate, Reset)

	sb.WriteString(DNGray + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(DNGray + "║" + Reset + "\n")

	sb.WriteString(DNGray + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *DeathNoteTheme) generateDNBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(DNDark + "〖" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DNDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DNDark + "〗" + Reset)
	return bar.String()
}
