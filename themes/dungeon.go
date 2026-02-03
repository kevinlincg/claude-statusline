package themes

import (
	"fmt"
	"strings"
)

// DungeonTheme 地牢火把風格
type DungeonTheme struct{}

func init() {
	RegisterTheme(&DungeonTheme{})
}

func (t *DungeonTheme) Name() string {
	return "dungeon"
}

func (t *DungeonTheme) Description() string {
	return "地牢：石牆火把照明，黑暗冒險氛圍"
}

const (
	DunStone     = "\033[38;2;105;105;105m"
	DunDarkStone = "\033[38;2;64;64;64m"
	DunTorch     = "\033[38;2;255;147;41m"
	DunFlame     = "\033[38;2;255;100;0m"
	DunGold      = "\033[38;2;255;215;0m"
	DunRed       = "\033[38;2;178;34;34m"
	DunGreen     = "\033[38;2;34;139;34m"
	DunBlue      = "\033[38;2;70;130;180m"
	DunPurple    = "\033[38;2;138;43;226m"
	DunBone      = "\033[38;2;255;250;240m"
	DunShadow    = "\033[38;2;25;25;25m"
	DunMoss      = "\033[38;2;85;107;47m"
)

func (t *DungeonTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Stone wall top with torches
	sb.WriteString(DunDarkStone + "▓▓▓" + DunTorch + "╔" + DunFlame + "҈" + DunTorch + "╗" + DunDarkStone + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + DunTorch + "╔" + DunFlame + "҈" + DunTorch + "╗" + DunDarkStone + "▓▓▓" + Reset + "\n")

	// Chamber name
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	chamberName := "The Dark Chamber"
	if data.ModelType == "Opus" {
		chamberName = "The Arcane Sanctum"
	} else if data.ModelType == "Haiku" {
		chamberName = "The Monk's Cell"
	}

	line1 := fmt.Sprintf("%s▓%s%s║%s%s %s%s%s%s  %s~ %s ~%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		DunTorch, chamberName, Reset)
	if data.UpdateAvailable {
		line1 += DunGold + " ★" + Reset
	}
	line1 += fmt.Sprintf("  %s%s%s", DunStone, data.Version, Reset)
	sb.WriteString(PadRight(line1, 77))
	sb.WriteString(DunTorch + "║" + DunDarkStone + "▓" + Reset + "\n")

	// Quest scroll
	line2 := fmt.Sprintf("%s▓%s%s║%s %s📜 %s%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunBone, ShortenPath(data.ProjectPath, 30), Reset)
	if data.GitBranch != "" {
		line2 += fmt.Sprintf("  %s⚔%s%s", DunMoss, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line2 += fmt.Sprintf(" %s+%d%s", DunGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line2 += fmt.Sprintf(" %s*%d%s", DunRed, data.GitDirty, Reset)
		}
	}
	sb.WriteString(PadRight(line2, 77))
	sb.WriteString(DunTorch + "║" + DunDarkStone + "▓" + Reset + "\n")

	// Stone separator
	sb.WriteString(DunDarkStone + "▓" + DunTorch + "╠" + DunStone + "══════════════════════════════════════════════════════════════════════════" + DunTorch + "╣" + DunDarkStone + "▓" + Reset + "\n")

	// Stats as dungeon items
	line3 := fmt.Sprintf("%s▓%s%s║%s %s🗡%s %s  %s🛡%s %d  %s⏳%s %s  %s💀%s %s  %s💎%s %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, FormatTokens(data.TokenCount),
		DunBlue, Reset, data.MessageCount,
		DunStone, Reset, data.SessionTime,
		DunPurple, Reset, FormatCostShort(data.BurnRate),
		DunGold, Reset, FormatCostShort(data.DayCost))
	sb.WriteString(PadRight(line3, 77))
	sb.WriteString(DunTorch + "║" + DunDarkStone + "▓" + Reset + "\n")

	// Health/Mana pools
	hp := 100 - data.ContextPercent
	hpColor := DunGreen
	if hp <= 20 {
		hpColor = DunRed
	} else if hp <= 50 {
		hpColor = DunTorch
	}

	hpBar := t.generateDungeonBar(hp, 15, hpColor)
	mpBar := t.generateDungeonBar(100-data.API5hrPercent, 12, DunBlue)
	xpBar := t.generateDungeonBar(100-data.API7dayPercent, 12, DunPurple)

	line4 := fmt.Sprintf("%s▓%s%s║%s %s❤%s%s%s%d%s  %s✦%s%s%s%d%s  %s⚡%s%s%s%d%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, hpBar, hpColor, hp, Reset,
		DunBlue, Reset, mpBar, DunBlue, 100-data.API5hrPercent, Reset,
		DunPurple, Reset, xpBar, DunPurple, 100-data.API7dayPercent, Reset)
	sb.WriteString(PadRight(line4, 77))
	sb.WriteString(DunTorch + "║" + DunDarkStone + "▓" + Reset + "\n")

	// Treasure info
	line5 := fmt.Sprintf("%s▓%s%s║%s %s💰%s %s ses  %s⚗%s %d%% hit  %s⌛%s %s  %s⌛%s %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunGold, Reset, FormatCostShort(data.SessionCost),
		DunGreen, Reset, data.CacheHitRate,
		DunStone, Reset, data.API5hrTimeLeft,
		DunStone, Reset, data.API7dayTimeLeft)
	sb.WriteString(PadRight(line5, 77))
	sb.WriteString(DunTorch + "║" + DunDarkStone + "▓" + Reset + "\n")

	// Stone wall bottom with torches
	sb.WriteString(DunDarkStone + "▓▓▓" + DunTorch + "╚" + DunFlame + "҈" + DunTorch + "╝" + DunDarkStone + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + DunTorch + "╚" + DunFlame + "҈" + DunTorch + "╝" + DunDarkStone + "▓▓▓" + Reset + "\n")

	return sb.String()
}

func (t *DungeonTheme) generateDungeonBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(DunShadow + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▰", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DunShadow)
		bar.WriteString(strings.Repeat("▱", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DunShadow + "〕" + Reset)
	return bar.String()
}
