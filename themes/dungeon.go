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

	// Fixed width: 78 characters
	width := 78

	// Stone wall top with torches
	torchTop := DunDarkStone + "▓▓▓" + DunTorch + "╔" + DunFlame + "*" + DunTorch + "╗" + DunDarkStone + strings.Repeat("▓", width-12) + DunTorch + "╔" + DunFlame + "*" + DunTorch + "╗" + DunDarkStone + "▓▓▓" + Reset
	sb.WriteString(torchTop + "\n")

	// Chamber name
	modelColor, _ := GetModelConfig(data.ModelType)
	chamberName := "The Dark Chamber"
	if data.ModelType == "Opus" {
		chamberName = "The Arcane Sanctum"
	} else if data.ModelType == "Haiku" {
		chamberName = "The Monk's Cell"
	}

	update := ""
	if data.UpdateAvailable {
		update = DunGold + " *" + Reset
	}

	line1 := fmt.Sprintf("%s▓%s%s|%s %s%s%s  %s~ %s ~%s%s  %s%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		modelColor, data.ModelName, Reset,
		DunTorch, chamberName, Reset, update,
		DunStone, data.Version)
	sb.WriteString(padLine(line1, width, DunTorch+"|"+DunDarkStone+"▓"+Reset))

	// Quest scroll
	gitStr := ""
	if data.GitBranch != "" {
		gitStr = fmt.Sprintf("  %s^%s%s", DunMoss, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitStr += fmt.Sprintf(" %s+%d%s", DunGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitStr += fmt.Sprintf(" %s*%d%s", DunRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("%s▓%s%s|%s %sScroll:%s %s%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunBone, Reset, ShortenPath(data.ProjectPath, 30), gitStr)
	sb.WriteString(padLine(line2, width, DunTorch+"|"+DunDarkStone+"▓"+Reset))

	// Stone separator
	sep := DunDarkStone + "▓" + DunTorch + "+" + DunStone + strings.Repeat("=", width-4) + DunTorch + "+" + DunDarkStone + "▓" + Reset
	sb.WriteString(sep + "\n")

	// Stats as dungeon items (no emojis)
	line3 := fmt.Sprintf("%s▓%s%s|%s %sSwd%s %s  %sShd%s %d  %sTime%s %s  %sSkul%s %s  %sGem%s %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, FormatTokens(data.TokenCount),
		DunBlue, Reset, data.MessageCount,
		DunStone, Reset, data.SessionTime,
		DunPurple, Reset, FormatCostShort(data.BurnRate),
		DunGold, Reset, FormatCostShort(data.DayCost))
	sb.WriteString(padLine(line3, width, DunTorch+"|"+DunDarkStone+"▓"+Reset))

	// Health/Mana pools
	hp := 100 - data.ContextPercent
	hpColor := DunGreen
	if hp <= 20 {
		hpColor = DunRed
	} else if hp <= 50 {
		hpColor = DunTorch
	}

	hpBar := t.generateDungeonBar(hp, 12, hpColor)
	mpBar := t.generateDungeonBar(100-data.API5hrPercent, 10, DunBlue)
	xpBar := t.generateDungeonBar(100-data.API7dayPercent, 10, DunPurple)

	line4 := fmt.Sprintf("%s▓%s%s|%s %sHP%s%s%s%d%s  %sMP%s%s%s%d%s  %sXP%s%s%s%d%s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunRed, Reset, hpBar, hpColor, hp, Reset,
		DunBlue, Reset, mpBar, DunBlue, 100-data.API5hrPercent, Reset,
		DunPurple, Reset, xpBar, DunPurple, 100-data.API7dayPercent, Reset)
	sb.WriteString(padLine(line4, width, DunTorch+"|"+DunDarkStone+"▓"+Reset))

	// Treasure info
	line5 := fmt.Sprintf("%s▓%s%s|%s %sGold%s %s ses  %sPotn%s %d%% hit  %sLeft%s %s / %s",
		DunDarkStone, Reset,
		DunTorch, Reset,
		DunGold, Reset, FormatCostShort(data.SessionCost),
		DunGreen, Reset, data.CacheHitRate,
		DunStone, Reset, data.API5hrTimeLeft, data.API7dayTimeLeft)
	sb.WriteString(padLine(line5, width, DunTorch+"|"+DunDarkStone+"▓"+Reset))

	// Stone wall bottom with torches
	torchBot := DunDarkStone + "▓▓▓" + DunTorch + "╚" + DunFlame + "*" + DunTorch + "╝" + DunDarkStone + strings.Repeat("▓", width-12) + DunTorch + "╚" + DunFlame + "*" + DunTorch + "╝" + DunDarkStone + "▓▓▓" + Reset
	sb.WriteString(torchBot + "\n")

	return sb.String()
}

// padLine pads a line to target visible width and adds suffix
func padLine(line string, targetWidth int, suffix string) string {
	visible := visibleLen(line)
	padding := targetWidth - visible - visibleLen(suffix)
	if padding < 0 {
		padding = 0
	}
	return line + strings.Repeat(" ", padding) + suffix + "\n"
}

// visibleLen calculates visible character count (excluding ANSI codes)
func visibleLen(s string) int {
	inEscape := false
	count := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
		} else {
			count++
		}
	}
	return count
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
	bar.WriteString(DunShadow + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("#", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(DunShadow)
		bar.WriteString(strings.Repeat("-", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(DunShadow + "]" + Reset)
	return bar.String()
}
