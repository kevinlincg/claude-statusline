package themes

import (
	"fmt"
	"strings"
)

// YokaiTheme Yokai mystical Japanese spirits style
type YokaiTheme struct{}

func init() {
	RegisterTheme(&YokaiTheme{})
}

func (t *YokaiTheme) Name() string {
	return "yokai"
}

func (t *YokaiTheme) Description() string {
	return "Yokai: Mystical Japanese spirits scroll style"
}

const (
	YKIPurple = "\033[38;2;100;50;150m"
	YKIBlue   = "\033[38;2;70;100;150m"
	YKIGreen  = "\033[38;2;100;150;100m"
	YKIRed    = "\033[38;2;150;50;50m"
	YKIGold   = "\033[38;2;180;150;80m"
	YKIWhite  = "\033[38;2;230;230;220m"
	YKIDark   = "\033[38;2;40;40;50m"
)

func (t *YokaiTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Mystical scroll opening
	sb.WriteString("\n")
	sb.WriteString("           " + YKIPurple + "｡" + YKIBlue + "ﾟ" + YKIGreen + ":" + YKIPurple + "｡" + YKIBlue + "ﾟ" + YKIGreen + ":" + Reset + "  " + YKIDark + "〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜" + Reset + "  " + YKIGreen + ":" + YKIBlue + "ﾟ" + YKIPurple + "｡" + YKIGreen + ":" + YKIBlue + "ﾟ" + YKIPurple + "｡" + Reset + "\n")
	sb.WriteString("\n")

	// Title in mystical style
	sb.WriteString("                    " + YKIPurple + "◈" + Reset + "  " + YKIWhite + "妖    怪    百    鬼    夜    行" + Reset + "  " + YKIPurple + "◈" + Reset + "\n")
	sb.WriteString("                              " + YKIGold + "～ YOKAI ～" + Reset + "\n")
	sb.WriteString("\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	spirit := "Kitsune"
	if data.ModelType == "Opus" {
		spirit = "Oni"
	} else if data.ModelType == "Haiku" {
		spirit = "Tanuki"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s〖覚醒〗%s", YKIPurple, Reset)
	}

	line1 := fmt.Sprintf("          %s霊:%s %s%s%s    %s形:%s %s%s%s    %s%s%s%s",
		YKIPurple, Reset, modelColor, modelIcon, data.ModelName,
		YKIBlue, Reset, YKIGold, spirit, Reset,
		YKIDark, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🌙%s%s", YKIBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", YKIGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", YKIRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("          %s界:%s %s%s",
		YKIGreen, Reset, ShortenPath(data.ProjectPath, 42), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString("\n")
	sb.WriteString("        " + YKIDark + "〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰" + Reset + "\n")
	sb.WriteString("\n")

	// Spirit powers
	yokiColor := YKIPurple
	if data.ContextPercent > 75 {
		yokiColor = YKIRed
	}

	line3 := fmt.Sprintf("              %s妖気%s    %s  %s%3d%%%s",
		YKIPurple, Reset, t.generateYKIBar(data.ContextPercent, 14, yokiColor), yokiColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("              %s霊力%s    %s  %s%3d%%%s  %s%s%s",
		YKIBlue, Reset, t.generateYKIBar(100-data.API5hrPercent, 14, YKIBlue),
		YKIBlue, 100-data.API5hrPercent, Reset, YKIDark, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("              %s呪力%s    %s  %s%3d%%%s  %s%s%s",
		YKIRed, Reset, t.generateYKIBar(data.API7dayPercent, 14, YKIRed),
		YKIRed, data.API7dayPercent, Reset, YKIDark, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString("\n")
	sb.WriteString("        " + YKIDark + "〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰〰" + Reset + "\n")
	sb.WriteString("\n")

	line6 := fmt.Sprintf("          %s%s%s 魂  %s%s%s  %s%d%s 術  %s%s%s 金  %s%d%%%s",
		YKIWhite, FormatTokens(data.TokenCount), Reset,
		YKIDark, data.SessionTime, Reset,
		YKIGreen, data.MessageCount, Reset,
		YKIGold, FormatCost(data.SessionCost), Reset,
		YKIPurple, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString("\n")
	sb.WriteString("           " + YKIPurple + "｡" + YKIBlue + "ﾟ" + YKIGreen + ":" + YKIPurple + "｡" + YKIBlue + "ﾟ" + YKIGreen + ":" + Reset + "  " + YKIDark + "〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜〜" + Reset + "  " + YKIGreen + ":" + YKIBlue + "ﾟ" + YKIPurple + "｡" + YKIGreen + ":" + YKIBlue + "ﾟ" + YKIPurple + "｡" + Reset + "\n")
	sb.WriteString("\n")

	return sb.String()
}

func (t *YokaiTheme) generateYKIBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(YKIDark + "〖" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("◉", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(YKIDark)
		bar.WriteString(strings.Repeat("◎", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(YKIDark + "〗" + Reset)
	return bar.String()
}
