package themes

import (
	"fmt"
	"strings"
)

// LaputaTheme Castle in the Sky style
type LaputaTheme struct{}

func init() {
	RegisterTheme(&LaputaTheme{})
}

func (t *LaputaTheme) Name() string {
	return "laputa"
}

func (t *LaputaTheme) Description() string {
	return "Laputa: Castle in the Sky flying stone style"
}

const (
	LPBlue    = "\033[38;2;100;149;237m"
	LPCyan    = "\033[38;2;0;200;200m"
	LPGreen   = "\033[38;2;144;238;144m"
	LPGold    = "\033[38;2;218;165;32m"
	LPWhite   = "\033[38;2;240;248;255m"
	LPGray    = "\033[38;2;119;136;153m"
	LPDark    = "\033[38;2;47;79;79m"
)

func (t *LaputaTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(LPBlue + "╔═════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(LPBlue + "║" + Reset + "  " + LPCyan + "◈" + LPWhite + " LAPUTA " + LPCyan + "◈" + Reset + "   " + LPBlue + "天空の城ラピュタ" + Reset + "   " + LPCyan + "Flying Stone Active" + Reset + "                   " + LPBlue + "║" + Reset + "\n")
	sb.WriteString(LPBlue + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	role := "Pazu"
	if data.ModelType == "Opus" {
		role = "Sheeta"
	} else if data.ModelType == "Haiku" {
		role = "Dola"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✧Levistone%s", LPCyan, Reset)
	}

	line1 := fmt.Sprintf("  %sPilot:%s %s%s%s  %sRole:%s %s%s%s  %s%s%s%s",
		LPCyan, Reset, modelColor, modelIcon, data.ModelName,
		LPGray, Reset, LPGold, role, Reset,
		LPGray, data.Version, Reset, update)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line1, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s☁%s%s", LPWhite, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", LPGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", LPGold, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sDestination:%s %s%s",
		LPGreen, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line2, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	sb.WriteString(LPBlue + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	stoneColor := LPCyan
	if data.ContextPercent > 75 {
		stoneColor = LPGold
	}

	line3 := fmt.Sprintf("  %sStone Power%s  %s  %s%3d%%%s",
		LPCyan, Reset, t.generateLPBar(data.ContextPercent, 18, stoneColor), stoneColor, data.ContextPercent, Reset)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line3, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sAltitude%s     %s  %s%3d%%%s  %s%s%s",
		LPBlue, Reset, t.generateLPBar(100-data.API5hrPercent, 18, LPBlue),
		LPBlue, 100-data.API5hrPercent, Reset, LPGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line4, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sRobots%s       %s  %s%3d%%%s  %s%s%s",
		LPGreen, Reset, t.generateLPBar(100-data.API7dayPercent, 18, LPGreen),
		LPGreen, 100-data.API7dayPercent, Reset, LPGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line5, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	sb.WriteString(LPBlue + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sData:%s %s%s%s  %sTime:%s %s  %sFlights:%s %s%d%s  %sTreasure:%s %s%s%s",
		LPWhite, Reset, LPWhite, FormatTokens(data.TokenCount), Reset,
		LPGray, Reset, data.SessionTime,
		LPGray, Reset, LPCyan, data.MessageCount, Reset,
		LPGold, Reset, LPGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line6, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sSync:%s %s%d%%%s",
		LPBlue, Reset, LPBlue, FormatCost(data.DayCost), Reset,
		LPGold, Reset, LPGold, FormatCost(data.BurnRate), Reset,
		LPCyan, Reset, LPCyan, data.CacheHitRate, Reset)

	sb.WriteString(LPBlue + "║" + Reset)
	sb.WriteString(PadRight(line7, 87))
	sb.WriteString(LPBlue + "║" + Reset + "\n")

	sb.WriteString(LPBlue + "╚═════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *LaputaTheme) generateLPBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(LPDark + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◆", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(LPDark)
		bar.WriteString(strings.Repeat("◇", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(LPDark + "〕" + Reset)
	return bar.String()
}
