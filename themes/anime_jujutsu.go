package themes

import (
	"fmt"
	"strings"
)

// JujutsuTheme Jujutsu Kaisen cursed energy style
type JujutsuTheme struct{}

func init() {
	RegisterTheme(&JujutsuTheme{})
}

func (t *JujutsuTheme) Name() string {
	return "jujutsu"
}

func (t *JujutsuTheme) Description() string {
	return "Jujutsu Kaisen: Cursed energy and domain expansion"
}

const (
	JJKPurple    = "\033[38;2;128;0;128m"
	JJKBlue      = "\033[38;2;0;100;255m"
	JJKRed       = "\033[38;2;200;0;0m"
	JJKBlack     = "\033[38;2;20;20;20m"
	JJKWhite     = "\033[38;2;240;240;240m"
	JJKPink      = "\033[38;2;255;100;150m"
	JJKCyan      = "\033[38;2;0;200;200m"
	JJKGold      = "\033[38;2;255;200;0m"
	JJKGray      = "\033[38;2;80;80;80m"
	JJKDark      = "\033[38;2;40;20;60m"
)

func (t *JujutsuTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Domain expansion style border
	sb.WriteString(JJKPurple + "╔═══════════════════════════════════════════════════════════════════════════════════════╗" + Reset + "\n")

	// Sorcerer info
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	grade := "Grade 2"
	if data.ModelType == "Opus" {
		grade = "Special Grade"
	} else if data.ModelType == "Haiku" {
		grade = "Grade 4"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s領域展開%s", JJKRed, Reset)
	}

	line1 := fmt.Sprintf(" %s呪術%s %s%s%s  %sGrade:%s %s%s%s  %s%s%s%s",
		JJKPurple, Reset,
		modelColor, modelIcon, data.ModelName,
		JJKGray, Reset, JJKGold, grade, Reset,
		JJKGray, data.Version, Reset, update)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line1, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	// Target curse
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s◈%s%s", JJKCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", JJKBlue, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s!%d%s", JJKRed, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf(" %sMission:%s %s%s",
		JJKRed, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line2, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	sb.WriteString(JJKPurple + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Cursed energy (context)
	ceColor := JJKBlue
	if data.ContextPercent > 75 {
		ceColor = JJKRed
	} else if data.ContextPercent > 50 {
		ceColor = JJKPurple
	}

	line3 := fmt.Sprintf(" %s呪力 Cursed Energy%s  %s  %s%3d%%%s",
		JJKBlue, Reset,
		t.generateJJKBar(data.ContextPercent, 16, ceColor),
		ceColor, data.ContextPercent, Reset)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line3, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	// Output limit (5hr)
	line4 := fmt.Sprintf(" %s出力 Output%s         %s  %s%3d%%%s  %s%s%s",
		JJKPink, Reset,
		t.generateJJKBar(100-data.API5hrPercent, 16, JJKPink),
		JJKPink, 100-data.API5hrPercent, Reset,
		JJKGray, data.API5hrTimeLeft, Reset)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line4, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	// Binding vow (7day)
	line5 := fmt.Sprintf(" %s縛り Binding Vow%s    %s  %s%3d%%%s  %s%s%s",
		JJKGold, Reset,
		t.generateJJKBar(100-data.API7dayPercent, 16, JJKGold),
		JJKGold, 100-data.API7dayPercent, Reset,
		JJKGray, data.API7dayTimeLeft, Reset)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line5, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	sb.WriteString(JJKPurple + "╠═══════════════════════════════════════════════════════════════════════════════════════╣" + Reset + "\n")

	// Stats
	line6 := fmt.Sprintf(" %sTechniques:%s %s%s%s  %sTime:%s %s  %sExorcisms:%s %s%d%s  %sBounty:%s %s%s%s",
		JJKCyan, Reset, JJKCyan, FormatTokens(data.TokenCount), Reset,
		JJKGray, Reset, data.SessionTime,
		JJKGray, Reset, JJKWhite, data.MessageCount, Reset,
		JJKGold, Reset, JJKGold, FormatCost(data.SessionCost), Reset)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line6, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	line7 := fmt.Sprintf(" %sDaily:%s %s%s%s  %sRate:%s %s%s/h%s  %sHit:%s %s%d%%%s",
		JJKPink, Reset, JJKPink, FormatCost(data.DayCost), Reset,
		JJKRed, Reset, JJKRed, FormatCost(data.BurnRate), Reset,
		JJKBlue, Reset, JJKBlue, data.CacheHitRate, Reset)

	sb.WriteString(JJKPurple + "║" + Reset)
	sb.WriteString(PadRight(line7, 89))
	sb.WriteString(JJKPurple + "║" + Reset + "\n")

	sb.WriteString(JJKPurple + "╚═══════════════════════════════════════════════════════════════════════════════════════╝" + Reset + "\n")

	return sb.String()
}

func (t *JujutsuTheme) generateJJKBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(JJKDark + "【" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▓", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(JJKDark)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(JJKDark + "】" + Reset)
	return bar.String()
}
