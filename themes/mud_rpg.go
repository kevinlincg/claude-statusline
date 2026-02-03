package themes

import (
	"fmt"
	"strings"
)

// MUDRPGTheme MUD RPG 角色狀態風格
type MUDRPGTheme struct{}

func init() {
	RegisterTheme(&MUDRPGTheme{})
}

func (t *MUDRPGTheme) Name() string {
	return "mud_rpg"
}

func (t *MUDRPGTheme) Description() string {
	return "MUD RPG：經典文字冒險遊戲角色狀態介面"
}

const (
	MUDGold     = "\033[38;2;255;215;0m"
	MUDRed      = "\033[38;2;255;80;80m"
	MUDBlue     = "\033[38;2;100;149;237m"
	MUDGreen    = "\033[38;2;50;205;50m"
	MUDCyan     = "\033[38;2;0;206;209m"
	MUDMagenta  = "\033[38;2;218;112;214m"
	MUDWhite    = "\033[38;2;245;245;245m"
	MUDGray     = "\033[38;2;128;128;128m"
	MUDDark     = "\033[38;2;64;64;64m"
	MUDBrown    = "\033[38;2;139;90;43m"
)

func (t *MUDRPGTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Header with character class style
	modelColor, _ := GetModelConfig(data.ModelType)
	className := "Artificer"
	if data.ModelType == "Opus" {
		className = "Archmage"
	} else if data.ModelType == "Haiku" {
		className = "Apprentice"
	}

	// Top border with ornate style
	sb.WriteString(MUDDark + "╔════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")

	// Character name and class
	update := ""
	if data.UpdateAvailable {
		update = MUDGold + " ★LEVEL UP★" + Reset
	}
	line1 := fmt.Sprintf("%s║%s  %s%s「%s」%s Lv.%s%s%s  %s◆ %s%s%s",
		MUDDark, Reset,
		modelColor, Bold, data.ModelName, Reset,
		MUDCyan, data.Version, Reset, update,
		MUDGray, ShortenPath(data.ProjectPath, 25), Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⚔%s%s", MUDBrown, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf("%s+%d%s", MUDGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf("%s~%d%s", MUDRed, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line1, 87))
	sb.WriteString(MUDDark + "║" + Reset + "\n")

	// Separator
	sb.WriteString(MUDDark + "╠════════════════════════════════════════╦═══════════════════════════════════════════╣" + Reset + "\n")

	// Stats row 1: HP (Context) / MP (5hr) / Class
	hpBar := t.generateHPBar(100-data.ContextPercent, 12) // Inverted: more context used = less HP
	mpBar := t.generateMPBar(100-data.API5hrPercent, 10)

	hpColor := MUDGreen
	if data.ContextPercent >= 80 {
		hpColor = MUDRed
	} else if data.ContextPercent >= 60 {
		hpColor = MUDGold
	}

	line2 := fmt.Sprintf("%s║%s  %sHP%s%s%s%3d%%%s  %sMP%s%s%s%3d%%%s  %s%s%s",
		MUDDark, Reset,
		MUDRed, Reset, hpBar, hpColor, 100-data.ContextPercent, Reset,
		MUDBlue, Reset, mpBar, MUDBlue, 100-data.API5hrPercent, Reset,
		MUDMagenta, className, Reset)
	sb.WriteString(PadRight(line2, 44))

	// Right column: XP and Gold
	xpBar := t.generateXPBar(data.API7dayPercent, 10)
	line2r := fmt.Sprintf("%s║%s  %sXP%s%s%s%3d%%%s  %sGP%s%s%s",
		MUDDark, Reset,
		MUDCyan, Reset, xpBar, MUDCyan, 100-data.API7dayPercent, Reset,
		MUDGold, MUDGold, FormatCostShort(data.DayCost), Reset)
	sb.WriteString(PadRight(line2r, 44))
	sb.WriteString(MUDDark + "║" + Reset + "\n")

	// Stats row 2: Combat stats
	line3 := fmt.Sprintf("%s║%s  %sATK%s %s  %sDEF%s %s  %sSPD%s %s  %sLUK%s %d%%",
		MUDDark, Reset,
		MUDRed, MUDWhite, FormatTokens(data.TokenCount),
		MUDBlue, MUDWhite, fmt.Sprintf("%d", data.MessageCount),
		MUDGreen, MUDWhite, data.SessionTime,
		MUDMagenta, MUDWhite, data.CacheHitRate)
	sb.WriteString(PadRight(line3, 44))

	// Right column: Equipment
	line3r := fmt.Sprintf("%s║%s  %s⚔%s %s  %s🛡%s %s/h  %s⏱%s %s",
		MUDDark, Reset,
		MUDBrown, MUDGold, FormatCostShort(data.SessionCost),
		MUDBrown, MUDRed, FormatCostShort(data.BurnRate),
		MUDGray, MUDGray, data.API5hrTimeLeft)
	sb.WriteString(PadRight(line3r, 44))
	sb.WriteString(MUDDark + "║" + Reset + "\n")

	// Bottom border
	sb.WriteString(MUDDark + "╚════════════════════════════════════════╩═══════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *MUDRPGTheme) generateHPBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MUDDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(MUDRed)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MUDDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MUDDark + "]" + Reset)
	return bar.String()
}

func (t *MUDRPGTheme) generateMPBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MUDDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(MUDBlue)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MUDDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MUDDark + "]" + Reset)
	return bar.String()
}

func (t *MUDRPGTheme) generateXPBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := (100 - percent) * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MUDDark + "[" + Reset)
	if filled > 0 {
		bar.WriteString(MUDCyan)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MUDDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MUDDark + "]" + Reset)
	return bar.String()
}
