package themes

import (
	"fmt"
	"strings"
)

// MahouTheme Magical Girl transformation style
type MahouTheme struct{}

func init() {
	RegisterTheme(&MahouTheme{})
}

func (t *MahouTheme) Name() string {
	return "mahou"
}

func (t *MahouTheme) Description() string {
	return "Mahou Shoujo: Magical girl transformation sparkle style"
}

const (
	MHPink    = "\033[38;2;255;105;180m"
	MHPurple  = "\033[38;2;186;85;211m"
	MHYellow  = "\033[38;2;255;215;0m"
	MHCyan    = "\033[38;2;0;255;255m"
	MHWhite   = "\033[38;2;255;250;250m"
	MHLavender = "\033[38;2;230;190;255m"
)

func (t *MahouTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Sparkly header
	sb.WriteString("    " + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + MHPink + "･ﾟ✧" + MHPurple + "･ﾟ" + MHYellow + "☆" + Reset + "\n")
	sb.WriteString("\n")
	sb.WriteString("           " + MHPink + "✨" + MHWhite + " 　Ｍ　Ａ　Ｈ　Ｏ　Ｕ　　Ｓ　Ｈ　Ｏ　Ｕ　Ｊ　Ｏ　 " + MHPink + "✨" + Reset + "\n")
	sb.WriteString("                          " + MHLavender + "～ 魔法少女 ～" + Reset + "\n")
	sb.WriteString("\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	magicalGirl := "Familiar"
	if data.ModelType == "Opus" {
		magicalGirl = "Guardian"
	} else if data.ModelType == "Haiku" {
		magicalGirl = "Mascot"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s✧･ﾟNew Power!･ﾟ✧%s", MHYellow, Reset)
	}

	line1 := fmt.Sprintf("      %s♡ Magical Girl:%s %s%s%s    %s♡ Form:%s %s%s%s  %s%s%s%s",
		MHPink, Reset, modelColor, modelIcon, data.ModelName,
		MHPurple, Reset, MHCyan, magicalGirl, Reset,
		MHLavender, data.Version, Reset, update)
	sb.WriteString(line1 + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s⋆%s%s", MHYellow, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", MHPink, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", MHPurple, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("      %s♡ Quest:%s %s%s",
		MHCyan, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)
	sb.WriteString(line2 + "\n")

	sb.WriteString("\n")
	sb.WriteString("    " + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + MHPurple + "･" + MHCyan + "｡" + MHPink + "･" + MHYellow + "｡" + Reset + "\n")
	sb.WriteString("\n")

	sparkleColor := MHPink
	if data.ContextPercent > 75 {
		sparkleColor = MHPurple
	}

	line3 := fmt.Sprintf("        %s✧ Sparkle Power%s  %s  %s%3d%%%s",
		MHPink, Reset, t.generateMHBar(data.ContextPercent, 16, sparkleColor), sparkleColor, data.ContextPercent, Reset)
	sb.WriteString(line3 + "\n")

	line4 := fmt.Sprintf("        %s✧ Love Energy%s    %s  %s%3d%%%s  %s%s%s",
		MHPurple, Reset, t.generateMHBar(100-data.API5hrPercent, 16, MHPurple),
		MHPurple, 100-data.API5hrPercent, Reset, MHLavender, data.API5hrTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	line5 := fmt.Sprintf("        %s✧ Hope Crystal%s   %s  %s%3d%%%s  %s%s%s",
		MHCyan, Reset, t.generateMHBar(100-data.API7dayPercent, 16, MHCyan),
		MHCyan, 100-data.API7dayPercent, Reset, MHLavender, data.API7dayTimeLeft, Reset)
	sb.WriteString(line5 + "\n")

	sb.WriteString("\n")
	sb.WriteString("    " + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + MHCyan + "｡" + MHPurple + "･" + MHPink + "｡" + MHYellow + "･" + Reset + "\n")
	sb.WriteString("\n")

	line6 := fmt.Sprintf("      %s%s%s stardust  %s%s%s  %s%d%s spells  %s$%s%s  %s%d%%%s shine",
		MHWhite, FormatTokens(data.TokenCount), Reset,
		MHLavender, data.SessionTime, Reset,
		MHPink, data.MessageCount, Reset,
		MHYellow, FormatCost(data.SessionCost), Reset,
		MHCyan, data.CacheHitRate, Reset)
	sb.WriteString(line6 + "\n")

	sb.WriteString("\n")
	sb.WriteString("    " + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + MHPurple + "･ﾟ✧" + MHPink + "･ﾟ" + MHYellow + "☆" + Reset + "\n")

	return sb.String()
}

func (t *MahouTheme) generateMHBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(MHLavender + "〔" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("♥", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(MHLavender)
		bar.WriteString(strings.Repeat("♡", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(MHLavender + "〕" + Reset)
	return bar.String()
}
