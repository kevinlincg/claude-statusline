package themes

import (
	"fmt"
	"strings"
)

// SchoolTheme Anime school blackboard/notebook style
type SchoolTheme struct{}

func init() {
	RegisterTheme(&SchoolTheme{})
}

func (t *SchoolTheme) Name() string {
	return "school"
}

func (t *SchoolTheme) Description() string {
	return "School: Anime school blackboard and notebook style"
}

const (
	SCHGreen  = "\033[38;2;50;100;50m"
	SCHWhite  = "\033[38;2;255;255;255m"
	SCHYellow = "\033[38;2;255;255;150m"
	SCHPink   = "\033[38;2;255;200;200m"
	SCHBlue   = "\033[38;2;150;200;255m"
	SCHChalk  = "\033[38;2;240;240;230m"
	SCHWood   = "\033[38;2;139;90;43m"
)

func (t *SchoolTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Blackboard frame
	sb.WriteString(SCHWood + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██████████████████████████████████████████████████████████████████████████████████" + SCHWood + "▓" + Reset + "\n")

	// Blackboard title (chalk writing)
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "          " + SCHChalk + "┌─────────────────────────────────────────────┐" + Reset + "             " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "          " + SCHChalk + "│" + Reset + "     " + SCHYellow + "★" + SCHWhite + " 学園アニメ " + SCHYellow + "★" + Reset + "  " + SCHChalk + "SCHOOL ANIME" + Reset + "      " + SCHChalk + "│" + Reset + "             " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "          " + SCHChalk + "└─────────────────────────────────────────────┘" + Reset + "             " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	modelColor, modelIcon := GetModelConfig(data.ModelType)
	student := "Class Rep"
	if data.ModelType == "Opus" {
		student = "Sensei"
	} else if data.ModelType == "Haiku" {
		student = "Kouhai"
	}

	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s[補習!]%s", SCHYellow, Reset)
	}

	line1 := fmt.Sprintf("    %s✎ Student:%s %s%s%s  %s✎ Role:%s %s%s%s  %s%s%s%s",
		SCHChalk, Reset, modelColor, modelIcon, data.ModelName,
		SCHChalk, Reset, SCHYellow, student, Reset,
		SCHChalk, data.Version, Reset, update)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line1, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s📚%s%s", SCHBlue, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", SCHPink, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", SCHYellow, data.GitDirty, Reset)
		}
	}

	line2 := fmt.Sprintf("    %s✎ Class:%s %s%s",
		SCHChalk, Reset, ShortenPath(data.ProjectPath, 45), gitInfo)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line2, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "    " + SCHChalk + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset + "    " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	// Grades/Stats
	gradeColor := SCHBlue
	if data.ContextPercent > 75 {
		gradeColor = SCHYellow
	}

	line3 := fmt.Sprintf("        %sHomework%s    %s  %s%3d%%%s",
		SCHPink, Reset, t.generateSCHBar(data.ContextPercent, 16, gradeColor), gradeColor, data.ContextPercent, Reset)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line3, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	line4 := fmt.Sprintf("        %sAttendance%s  %s  %s%3d%%%s  %s%s%s",
		SCHBlue, Reset, t.generateSCHBar(100-data.API5hrPercent, 16, SCHBlue),
		SCHBlue, 100-data.API5hrPercent, Reset, SCHChalk, data.API5hrTimeLeft, Reset)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line4, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	line5 := fmt.Sprintf("        %sGrades%s      %s  %s%3d%%%s  %s%s%s",
		SCHYellow, Reset, t.generateSCHBar(100-data.API7dayPercent, 16, SCHYellow),
		SCHYellow, 100-data.API7dayPercent, Reset, SCHChalk, data.API7dayTimeLeft, Reset)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line5, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "    " + SCHChalk + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + Reset + "    " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	line6 := fmt.Sprintf("        %s%s%s notes  %s%s%s  %s%d%s periods  %s¥%s%s  %s%d%%%s",
		SCHWhite, FormatTokens(data.TokenCount), Reset,
		SCHChalk, data.SessionTime, Reset,
		SCHBlue, data.MessageCount, Reset,
		SCHYellow, FormatCost(data.SessionCost), Reset,
		SCHPink, data.CacheHitRate, Reset)

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset)
	sb.WriteString(PadRight(line6, 78))
	sb.WriteString(SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")

	sb.WriteString(SCHWood + "▓" + SCHGreen + "██" + Reset + "                                                                                " + SCHGreen + "██" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓" + SCHGreen + "██████████████████████████████████████████████████████████████████████████████████" + SCHWood + "▓" + Reset + "\n")
	sb.WriteString(SCHWood + "▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓" + Reset + "\n")

	return sb.String()
}

func (t *SchoolTheme) generateSCHBar(percent, width int, color string) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	var bar strings.Builder
	bar.WriteString(SCHChalk + "[" + Reset)
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(SCHChalk)
		bar.WriteString(strings.Repeat("░", empty))
		bar.WriteString(Reset)
	}
	bar.WriteString(SCHChalk + "]" + Reset)
	return bar.String()
}
