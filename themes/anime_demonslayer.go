package themes

import (
	"fmt"
	"strings"
)

// DemonSlayerTheme Demon Slayer breathing style
type DemonSlayerTheme struct{}

func init() {
	RegisterTheme(&DemonSlayerTheme{})
}

func (t *DemonSlayerTheme) Name() string {
	return "demonslayer"
}

func (t *DemonSlayerTheme) Description() string {
	return "Demon Slayer: Breathing techniques and Nichirin blade"
}

const (
	DSRed        = "\033[38;2;220;20;60m"
	DSBlue       = "\033[38;2;30;144;255m"
	DSYellow     = "\033[38;2;255;215;0m"
	DSGreen      = "\033[38;2;50;205;50m"
	DSPink       = "\033[38;2;255;182;193m"
	DSPurple     = "\033[38;2;148;0;211m"
	DSOrange     = "\033[38;2;255;140;0m"
	DSWhite      = "\033[38;2;255;255;255m"
	DSDark       = "\033[38;2;40;40;40m"
	DSGray       = "\033[38;2;100;100;100m"
)

func (t *DemonSlayerTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Nichirin blade inspired border
	sb.WriteString(DSRed + "═══════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	// Corps info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	breathStyle := "Water"
	breathColor := DSBlue
	if data.ModelType == "Opus" {
		breathStyle = "Sun"
		breathColor = DSYellow
	} else if data.ModelType == "Haiku" {
		breathStyle = "Flower"
		breathColor = DSPink
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[鬼殺隊]%s", DSRed, Reset)
	}

	line1 := fmt.Sprintf(" %s鬼滅%s %s%s%s  %sBreath:%s %s%s%s  %s%s%s%s",
		DSRed, Reset,
		modelColor, modelIcon, data.ModelName,
		DSGray, Reset, breathColor, breathStyle, Reset,
		DSGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	// Target demon (project)
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⚔%s%s", DSBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", DSGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", DSOrange, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf(" %sTarget:%s %s%s",
		DSPurple, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(DSGray + "─────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	// Breathing gauge (context)
	breathGaugeColor := breathColor
	if data.ContextPercent > 75 {
		breathGaugeColor = DSRed
	} else if data.ContextPercent > 50 {
		breathGaugeColor = DSOrange
	}

	line3 := fmt.Sprintf(" %s呼吸 Breath%s    %s  %s%3d%%%s",
		breathColor, Reset,
		t.generateDSBar(data.ContextPercent, 18, breathGaugeColor),
		breathGaugeColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	// Stamina (5hr)
	line4 := fmt.Sprintf(" %s体力 Stamina%s   %s  %s%3d%%%s  %s%s%s",
		DSGreen, Reset,
		t.generateDSBar(100-data.API5hrPercent, 18, DSGreen),
		DSGreen, 100-data.API5hrPercent, Reset,
		DSGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	// Focus (7day)
	line5 := fmt.Sprintf(" %s集中 Focus%s     %s  %s%3d%%%s  %s%s%s",
		DSPurple, Reset,
		t.generateDSBar(100-data.API7dayPercent, 18, DSPurple),
		DSPurple, 100-data.API7dayPercent, Reset,
		DSGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(DSGray + "─────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	// Combat stats
	line6 := fmt.Sprintf(" %sForms:%s %s%s%s  %sTime:%s %s  %sSlays:%s %s%d%s  %sReward:%s %s%s%s  %sDaily:%s %s%s%s",
		DSBlue, Reset, DSBlue, FormatTokens(data.TokenCount), Reset,
		DSGray, Reset, data.SessionTime,
		DSGray, Reset, DSWhite, data.MessageCount, Reset,
		DSYellow, Reset, DSYellow, FormatCost(data.SessionCost), Reset,
		DSOrange, Reset, DSOrange, FormatCost(data.DayCost), Reset)
	sb.WriteString(line6 + "\n")

	line7 := fmt.Sprintf(" %sRate:%s %s%s/h%s  %sPrecision:%s %s%d%%%s",
		DSRed, Reset, DSRed, FormatCost(data.BurnRate), Reset,
		DSGreen, Reset, DSGreen, data.CacheHitRate, Reset)
	sb.WriteString(line7 + "\n")

	sb.WriteString(DSRed + "═══════════════════════════════════════════════════════════════════════════════════════" + Reset + "\n")

	return sb.String()
}

func (t *DemonSlayerTheme) generateDSBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(DSDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◆", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DSDark)
		bar.WriteString(strings.Repeat("◇", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DSDark + "〉" + Reset)
	return bar.String()
}
