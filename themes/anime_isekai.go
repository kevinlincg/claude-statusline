package themes

import (
	"fmt"
	"strings"
)

// IsekaiTheme Generic isekai RPG status window style
type IsekaiTheme struct{}

func init() {
	RegisterTheme(&IsekaiTheme{})
}

func (t *IsekaiTheme) Name() string {
	return "isekai"
}

func (t *IsekaiTheme) Description() string {
	return "Isekai: RPG status window, HP/MP/EXP style"
}

const (
	IsekaiGold       = "\033[38;2;255;215;0m"
	IsekaiBlue       = "\033[38;2;100;149;237m"
	IsekaiGreen      = "\033[38;2;50;205;50m"
	IsekaiRed        = "\033[38;2;220;20;60m"
	IsekaiPurple     = "\033[38;2;148;0;211m"
	IsekaiWhite      = "\033[38;2;255;255;255m"
	IsekaiDark       = "\033[38;2;64;64;64m"
	IsekaiCyan       = "\033[38;2;0;255;255m"
	IsekaiBrown      = "\033[38;2;139;69;19m"
	IsekaiBgBlue     = "\033[48;2;20;30;60m"
)

func (t *IsekaiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// RPG Window frame
	sb.WriteString(IsekaiGold + "╔════════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")

	// Character name and class
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	className := "Sage"
	if data.ModelType == "Opus" {
		className = "Archmage"
	} else if data.ModelType == "Haiku" {
		className = "Apprentice"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[LEVEL UP!]%s", IsekaiCyan, Reset)
	}

	line1 := fmt.Sprintf(" %s%s%s %s%s  %sLv.%s%s  %sClass:%s %s%s%s%s",
		modelColor, modelIcon, data.ModelName, Reset,
		IsekaiDark, IsekaiWhite, data.Version, Reset,
		IsekaiGold, Reset, IsekaiPurple, className, Reset, update)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line1, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	sb.WriteString(IsekaiGold + "╠════════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Quest info
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf(" %s⚔%s%s", IsekaiCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", IsekaiGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s!%d%s", IsekaiRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf(" %sQuest:%s %s%s",
		IsekaiGold, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line2, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	sb.WriteString(IsekaiGold + "╠════════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// HP bar (context - inverted, more is bad)
	hpPercent := 100 - data.ContextPercent
	hpColor := IsekaiGreen
	if hpPercent < 25 {
		hpColor = IsekaiRed
	} else if hpPercent < 50 {
		hpColor = IsekaiGold
	}

	line3 := fmt.Sprintf(" %sHP%s  %s %s%3d/100%s",
		IsekaiRed, Reset,
		t.generateIsekaiBar(hpPercent, 25, hpColor, IsekaiRed),
		hpColor, hpPercent, Reset)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line3, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	// MP bar (5hr limit)
	mpPercent := 100 - data.API5hrPercent
	mpColor := IsekaiBlue
	if mpPercent < 25 {
		mpColor = IsekaiPurple
	}

	line4 := fmt.Sprintf(" %sMP%s  %s %s%3d/100%s  %sRegen:%s %s",
		IsekaiBlue, Reset,
		t.generateIsekaiBar(mpPercent, 25, mpColor, IsekaiBlue),
		mpColor, mpPercent, Reset,
		IsekaiDark, Reset, data.API5hrTimeLeft)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line4, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	// Stamina bar (7day limit)
	staminaPercent := 100 - data.API7dayPercent
	staminaColor := IsekaiGreen
	if staminaPercent < 25 {
		staminaColor = IsekaiRed
	}

	line5 := fmt.Sprintf(" %sSP%s  %s %s%3d/100%s  %sRegen:%s %s",
		IsekaiGreen, Reset,
		t.generateIsekaiBar(staminaPercent, 25, staminaColor, IsekaiGreen),
		staminaColor, staminaPercent, Reset,
		IsekaiDark, Reset, data.API7dayTimeLeft)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line5, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	sb.WriteString(IsekaiGold + "╠════════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Stats panel
	line6 := fmt.Sprintf(" %sEXP:%s %s%s%s  %sTime:%s %s  %sActions:%s %s%d%s",
		IsekaiPurple, Reset, IsekaiCyan, FormatTokens(data.TokenCount), Reset,
		IsekaiDark, Reset, data.SessionTime,
		IsekaiDark, Reset, IsekaiWhite, data.MessageCount, Reset)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line6, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	line7 := fmt.Sprintf(" %sGold:%s %s%s%s  %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sLuck:%s %s%d%%%s",
		IsekaiGold, Reset, IsekaiGold, FormatCost(data.SessionCost), Reset,
		IsekaiBrown, Reset, IsekaiBrown, FormatCost(data.DayCost), Reset,
		IsekaiRed, Reset, IsekaiRed, FormatCost(data.BurnRate), Reset,
		IsekaiCyan, Reset, IsekaiCyan, data.CacheHitRate, Reset)

	sb.WriteString(IsekaiGold + "║" + Reset)
	sb.WriteString(PadRight(line7, 90))
	sb.WriteString(IsekaiGold + "║" + Reset + "\n")

	sb.WriteString(IsekaiGold + "╚════════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *IsekaiTheme) generateIsekaiBar(percent, width int, fillColor, borderColor string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(borderColor + "【" + Reset)
	if filled > 0 {
		bar.WriteString(fillColor)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(IsekaiDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(borderColor + "】" + Reset)
	return bar.String()
}
