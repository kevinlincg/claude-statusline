package themes

import (
	"fmt"
	"strings"
)

// TokyoGhoulTheme Tokyo Ghoul kagune style
type TokyoGhoulTheme struct{}

func init() {
	RegisterTheme(&TokyoGhoulTheme{})
}

func (t *TokyoGhoulTheme) Name() string {
	return "tokyoghoul"
}

func (t *TokyoGhoulTheme) Description() string {
	return "Tokyo Ghoul: Ghoul kagune RC cell style"
}

const (
	TGRed    = "\033[38;2;139;0;0m"
	TGBlack  = "\033[38;2;20;20;20m"
	TGWhite  = "\033[38;2;240;240;240m"
	TGPurple = "\033[38;2;100;50;100m"
	TGGray   = "\033[38;2;80;80;80m"
	TGDark   = "\033[38;2;30;30;30m"
)

func (t *TokyoGhoulTheme) Render(data StatusData) string {
	var sb strings.Builder

	sb.WriteString(TGRed + "█████████████████████████████████████████████████████████████████████████████████████████" + Reset + "\n")
	sb.WriteString(TGBlack + "███" + Reset + " " + TGRed + "👁" + TGWhite + " TOKYO GHOUL " + TGRed + "👁" + Reset + "   " + TGPurple + "東京喰種" + Reset + "   " + TGRed + "// RC CELL ACTIVE //" + Reset + "               " + TGBlack + "███" + Reset + "\n")
	sb.WriteString(TGRed + "█████████████████████████████████████████████████████████████████████████████████████████" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	ghoul := "Hinami"
	if data.ModelType == "Opus" {
		ghoul = "Kaneki"
	} else if data.ModelType == "Haiku" {
		ghoul = "Touka"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[Awakening]%s", TGRed, Reset)
	}

	line1 := fmt.Sprintf("  %sGhoul:%s %s%s%s  %sKagune:%s %s%s%s  %s%s%s%s",
		TGRed, Reset, modelColor, modelIcon, data.ModelName,
		TGGray, Reset, TGPurple, ghoul, Reset,
		TGGray, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s☕%s%s", TGPurple, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", TGWhite, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", TGRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("  %sTerritory:%s %s%s",
		TGPurple, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString(TGRed + "─────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	rcColor := TGPurple
	if data.ContextPercent > 75 {
		rcColor = TGRed
	}

	line3 := fmt.Sprintf("  %sRC Cells%s   %s  %s%3d%%%s",
		TGRed, Reset, t.generateTGBar(data.ContextPercent, 18, rcColor), rcColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("  %sHunger%s     %s  %s%3d%%%s  %s%s%s",
		TGPurple, Reset, t.generateTGBar(data.API5hrPercent, 18, TGPurple),
		TGPurple, data.API5hrPercent, Reset, TGGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("  %sKagune%s     %s  %s%3d%%%s  %s%s%s",
		TGRed, Reset, t.generateTGBar(100-data.API7dayPercent, 18, TGRed),
		TGRed, 100-data.API7dayPercent, Reset, TGGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString(TGRed + "─────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	line6 := fmt.Sprintf("  %s%s%s cells  %s%s%s  %s%d%s hunts  %s$%s%s  %s$%s/day%s  %s%d%%%s rate",
		TGWhite, FormatTokens(data.TokenCount), Reset,
		TGGray, data.SessionTime, Reset,
		TGPurple, data.MessageCount, Reset,
		TGRed, FormatCost(data.SessionCost), Reset,
		TGPurple, FormatCost(data.DayCost), Reset,
		TGRed, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString(TGRed + "█████████████████████████████████████████████████████████████████████████████████████████" + Reset + "\n")

	return sb.String()
}

func (t *TokyoGhoulTheme) generateTGBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(TGDark + "〈" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(TGDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(TGDark + "〉" + Reset)
	return bar.String()
}
