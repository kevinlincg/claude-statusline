package themes

import (
	"fmt"
	"strings"
)

// HowlTheme Howl's Moving Castle style
type HowlTheme struct{}

func init() {
	RegisterTheme(&HowlTheme{})
}

func (t *HowlTheme) Name() string {
	return "howl"
}

func (t *HowlTheme) Description() string {
	return "Howl: Moving Castle steam magic style"
}

const (
	HowlCopper  = "\033[38;2;184;115;51m"
	HowlGold    = "\033[38;2;255;215;0m"
	HowlOrange  = "\033[38;2;255;140;0m"
	HowlBlue    = "\033[38;2;135;206;250m"
	HowlPurple  = "\033[38;2;147;112;219m"
	HowlWhite   = "\033[38;2;255;250;240m"
	HowlGray    = "\033[38;2;120;120;120m"
	HowlDark    = "\033[38;2;60;40;30m"
)

func (t *HowlTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(HowlCopper + "╔═════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")
	sb.WriteString(HowlCopper + "║" + Reset + "  " + HowlOrange + "🔥" + HowlWhite + " Moving Castle " + HowlOrange + "🔥" + Reset + "   " + HowlPurple + "ハウルの動く城" + Reset + "                                    " + HowlCopper + "║" + Reset + "\n")
	sb.WriteString(HowlCopper + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	resident := "Sophie"
	if data.ModelType == "Opus" {
		resident = "Howl"
	} else if data.ModelType == "Haiku" {
		resident = "Markl"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✨Magic!%s", HowlPurple, Reset)
	}

	line1 := fmt.Sprintf("  %sWizard:%s %s%s%s  %sName:%s %s%s%s  %s%s%s%s",
		HowlPurple, Reset, modelColor, modelIcon, data.ModelName,
		HowlGray, Reset, HowlBlue, resident, Reset,
		HowlGray, data.Version, Reset, update)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line1, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚙%s%s", HowlCopper, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", HowlGold, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", HowlOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sDoor:%s %s%s",
		HowlBlue, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line2, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	sb.WriteString(HowlCopper + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	fireColor := HowlOrange
	if data.ContextPercent > 75 {
		fireColor = HowlOrange
	}

	line3 := fmt.Sprintf("  %sCalcifer%s   %s  %s%3d%%%s",
		HowlOrange, Reset, t.generateHowlBar(data.ContextPercent, 18, fireColor), fireColor, data.ContextPercent, Reset)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line3, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	line4 := fmt.Sprintf("  %sMagic%s      %s  %s%3d%%%s  %s%s%s",
		HowlPurple, Reset, t.generateHowlBar(100-data.API5hrPercent, 18, HowlPurple),
		HowlPurple, 100-data.API5hrPercent, Reset, HowlGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line4, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	line5 := fmt.Sprintf("  %sSteam%s      %s  %s%3d%%%s  %s%s%s",
		HowlCopper, Reset, t.generateHowlBar(100-data.API7dayPercent, 18, HowlCopper),
		HowlCopper, 100-data.API7dayPercent, Reset, HowlGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line5, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	sb.WriteString(HowlCopper + "╠═════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	line6 := fmt.Sprintf("  %sSpells:%s %s%s%s  %sTime:%s %s  %sDeliveries:%s %s%d%s  %sGold:%s %s%s%s",
		HowlWhite, Reset, HowlWhite, FormatTokens(data.TokenCount), Reset,
		HowlGray, Reset, data.SessionTime,
		HowlGray, Reset, HowlBlue, data.MessageCount, Reset,
		HowlGold, Reset, HowlGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line6, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	line7 := fmt.Sprintf("  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sHeart:%s %s%d%%%s",
		HowlPurple, Reset, HowlPurple, FormatCost(data.DayCost), Reset,
		HowlOrange, Reset, HowlOrange, FormatCost(data.BurnRate), Reset,
		HowlBlue, Reset, HowlBlue, data.CacheHitRate, Reset)

	sb.WriteString(HowlCopper + "║" + Reset)
	sb.WriteString(PadRight(line7, 87))
	sb.WriteString(HowlCopper + "║" + Reset + "\n")

	sb.WriteString(HowlCopper + "╚═════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *HowlTheme) generateHowlBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(HowlDark + "⟨" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(HowlDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(HowlDark + "⟩" + Reset)
	return bar.String()
}
