package themes

import (
	"fmt"
	"strings"
)

// OnelineCleanTheme single-line with dot separators
type OnelineCleanTheme struct{}

func init() {
	RegisterTheme(&OnelineCleanTheme{})
}

func (t *OnelineCleanTheme) Name() string {
	return "oneline_clean"
}

func (t *OnelineCleanTheme) Description() string {
	return "Single-line clean: dot separators, colored text on dark background"
}

func (t *OnelineCleanTheme) Render(data StatusData) string {
	var sb strings.Builder

	sep := fmt.Sprintf(" %s·%s ", ColorDim, Reset)

	// Model + Version
	modelColor, _ := GetModelConfig(data.ModelType)
	sb.WriteString(fmt.Sprintf(" %s%s%s%s", modelColor, Bold, data.ModelName, Reset))
	sb.WriteString(fmt.Sprintf(" %s%s%s", ColorDim, data.Version, Reset))

	sb.WriteString(sep)

	// Path
	sb.WriteString(fmt.Sprintf("%s%s%s", ColorBlue, ShortenPath(data.ProjectPath, 20), Reset))

	// Git
	if data.GitBranch != "" {
		sb.WriteString(sep)
		sb.WriteString(fmt.Sprintf("%s%s%s", ColorGreen, data.GitBranch, Reset))
		if data.GitStaged > 0 {
			sb.WriteString(fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset))
		}
		if data.GitDirty > 0 {
			sb.WriteString(fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset))
		}
	}

	sb.WriteString(sep)

	// 5hr
	color5, _ := GetBarColor(data.API5hrPercent)
	sb.WriteString(fmt.Sprintf("%s5h%s %s%d%%%s", ColorDim, Reset, color5, data.API5hrPercent, Reset))
	if data.API5hrTimeLeft != "" {
		sb.WriteString(fmt.Sprintf(" %s%s%s", ColorDim, data.API5hrTimeLeft, Reset))
	}

	sb.WriteString(sep)

	// 7day
	color7, _ := GetBarColor(data.API7dayPercent)
	sb.WriteString(fmt.Sprintf("%s7d%s %s%d%%%s", ColorDim, Reset, color7, data.API7dayPercent, Reset))
	if data.API7dayTimeLeft != "" {
		sb.WriteString(fmt.Sprintf(" %s%s%s", ColorDim, data.API7dayTimeLeft, Reset))
	}

	sb.WriteString("\n")
	return sb.String()
}
