package themes

import (
	"fmt"
	"strings"
)

// OnelinePillsTheme single-line with rounded pill badges
type OnelinePillsTheme struct{}

func init() {
	RegisterTheme(&OnelinePillsTheme{})
}

func (t *OnelinePillsTheme) Name() string {
	return "oneline_pills"
}

func (t *OnelinePillsTheme) Description() string {
	return "Single-line pills: rounded badge segments with mini progress bars"
}

// Pill colors
const (
	PillDim    = "\033[38;2;100;100;100m"
	PillBorder = "\033[38;2;180;180;180m"
)

func (t *OnelinePillsTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Pill: Model
	modelColor, _ := GetModelConfig(data.ModelType)
	sb.WriteString(t.pill(
		fmt.Sprintf("%s%s%s %s%s%s", modelColor, Bold, data.ModelName, Reset+ColorDim, data.Version, Reset),
		PillBorder))

	sb.WriteString(" ")

	// Pill: Path
	sb.WriteString(t.pill(
		fmt.Sprintf("%s%s%s", ColorBlue, ShortenPath(data.ProjectPath, 20), Reset),
		PillBorder))

	// Pill: Git
	if data.GitBranch != "" {
		sb.WriteString(" ")
		gitContent := fmt.Sprintf("%s%s%s", ColorGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitContent += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitContent += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
		sb.WriteString(t.pill(gitContent, PillBorder))
	}

	sb.WriteString(" ")

	// Pill: 5hr with mini bar
	color5, _ := GetBarColor(data.API5hrPercent)
	bar5 := t.miniBar(data.API5hrPercent, 6, color5)
	time5 := ""
	if data.API5hrTimeLeft != "" {
		time5 = fmt.Sprintf(" %s%s%s", ColorDim, data.API5hrTimeLeft, Reset)
	}
	sb.WriteString(t.pill(
		fmt.Sprintf("%s5h%s %s %s%d%%%s%s", ColorDim, Reset, bar5, color5, data.API5hrPercent, Reset, time5),
		PillBorder))

	sb.WriteString(" ")

	// Pill: 7day with mini bar
	color7, _ := GetBarColor(data.API7dayPercent)
	bar7 := t.miniBar(data.API7dayPercent, 6, color7)
	time7 := ""
	if data.API7dayTimeLeft != "" {
		time7 = fmt.Sprintf(" %s%s%s", ColorDim, data.API7dayTimeLeft, Reset)
	}
	sb.WriteString(t.pill(
		fmt.Sprintf("%s7d%s %s %s%d%%%s%s", ColorDim, Reset, bar7, color7, data.API7dayPercent, Reset, time7),
		PillBorder))

	sb.WriteString("\n")
	return sb.String()
}

func (t *OnelinePillsTheme) pill(content, borderColor string) string {
	return fmt.Sprintf("%s(%s %s %s)%s", borderColor, Reset, content, borderColor, Reset)
}

func (t *OnelinePillsTheme) miniBar(percent, width int, color string) string {
	filled := percent * width / 100
	if filled > width {
		filled = width
	}
	empty := width - filled

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("▮", filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(PillDim)
		bar.WriteString(strings.Repeat("▯", empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
