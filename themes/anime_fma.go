package themes

import (
	"fmt"
	"strings"
)

// FMATheme Fullmetal Alchemist transmutation style
type FMATheme struct{}

func init() {
	RegisterTheme(&FMATheme{})
}

func (t *FMATheme) Name() string {
	return "fma"
}

func (t *FMATheme) Description() string {
	return "FMA: Fullmetal Alchemist transmutation circle style"
}

const (
	FMAGold    = "\033[38;2;218;165;32m"
	FMARed     = "\033[38;2;180;0;0m"
	FMABlue    = "\033[38;2;70;130;180m"
	FMAWhite   = "\033[38;2;240;240;240m"
	FMAGray    = "\033[38;2;120;120;120m"
	FMADark    = "\033[38;2;50;50;50m"
	FMAPurple  = "\033[38;2;128;0;128m"
)

func (t *FMATheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(FMAGold + "╔══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(FMAGold + "║" + Reset + "           " + FMAWhite + "☆ EQUIVALENT EXCHANGE ☆" + Reset + "   " + FMAGold + "「等価交換」" + Reset + "                            " + FMAGold + "║" + Reset + "\n")
	sb.WriteString(FMAGold + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	rank := "State Alchemist"
	if data.ModelType == "Opus" {
		rank = "Führer's Alchemist"
	} else if data.ModelType == "Haiku" {
		rank = "Apprentice"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[TRANSMUTE!]%s", FMAPurple, Reset)
	}

	line1 := fmt.Sprintf("  %sAlchemist:%s %s%s%s  %sRank:%s %s%s%s  %s%s%s%s",
		FMAGold, Reset, modelColor, modelIcon, data.ModelName,
		FMAGray, Reset, FMABlue, rank, Reset,
		FMAGray, data.Version, Reset, update)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line1, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚗%s%s", FMABlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", FMAGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", FMARed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sResearch:%s %s%s",
		FMARed, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line2, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	sb.WriteString(FMAGold + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	energyColor := FMAGold
	if data.ContextPercent > 75 {
		energyColor = FMARed
	}

	line3 := fmt.Sprintf("  %sAlchemic Energy%s  %s  %s%3d%%%s",
		FMAGold, Reset, t.generateFMABar(data.ContextPercent, 18, energyColor), energyColor, data.ContextPercent, Reset)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line3, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sPhysical%s         %s  %s%3d%%%s  %s%s%s",
		FMABlue, Reset, t.generateFMABar(100-data.API5hrPercent, 18, FMABlue),
		FMABlue, 100-data.API5hrPercent, Reset, FMAGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line4, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sSoul%s             %s  %s%3d%%%s  %s%s%s",
		FMAPurple, Reset, t.generateFMABar(100-data.API7dayPercent, 18, FMAPurple),
		FMAPurple, 100-data.API7dayPercent, Reset, FMAGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line5, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	sb.WriteString(FMAGold + "╠══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sTransmutations:%s %s%s%s  %sTime:%s %s  %sCircles:%s %s%d%s  %sCenz:%s %s%s%s",
		FMAWhite, Reset, FMAWhite, FormatTokens(data.TokenCount), Reset,
		FMAGray, Reset, data.SessionTime,
		FMAGray, Reset, FMAGold, data.MessageCount, Reset,
		FMAGold, Reset, FMAGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line6, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sEfficiency:%s %s%d%%%s",
		FMABlue, Reset, FMABlue, FormatCost(data.DayCost), Reset,
		FMARed, Reset, FMARed, FormatCost(data.BurnRate), Reset,
		FMAGold, Reset, FMAGold, data.CacheHitRate, Reset)

	sb.WriteString(FMAGold + "║" + Reset)
	sb.WriteString(PadRight(line7, 88))
	sb.WriteString(FMAGold + "║" + Reset + "\n")

	sb.WriteString(FMAGold + "╚══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *FMATheme) generateFMABar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(FMADark + "⟨" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◈", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(FMADark)
		bar.WriteString(strings.Repeat("◇", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(FMADark + "⟩" + Reset)
	return bar.String()
}
