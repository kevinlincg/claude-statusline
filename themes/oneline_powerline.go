package themes

import (
	"fmt"
	"strings"
)

// OnelinePowerlineTheme single-line with powerline arrow segments
type OnelinePowerlineTheme struct{}

func init() {
	RegisterTheme(&OnelinePowerlineTheme{})
}

func (t *OnelinePowerlineTheme) Name() string {
	return "oneline_powerline"
}

func (t *OnelinePowerlineTheme) Description() string {
	return "Single-line powerline: colored arrow segments like a shell prompt"
}

// Powerline segment background/foreground colors
const (
	// Model segment: dark gold/amber
	PLModelBg  = "\033[48;2;60;50;30m"
	PLModelFg  = "\033[38;2;60;50;30m"
	PLModelTxt = "\033[38;2;220;185;100m"

	// Path segment: dark blue
	PLPathBg  = "\033[48;2;30;40;65m"
	PLPathFg  = "\033[38;2;30;40;65m"
	PLPathTxt = "\033[38;2;130;170;230m"

	// Git segment: dark green
	PLGitBg  = "\033[48;2;30;55;35m"
	PLGitFg  = "\033[38;2;30;55;35m"
	PLGitTxt = "\033[38;2;130;200;140m"

	// 5hr segment: dark teal
	PL5hrBg  = "\033[48;2;25;50;55m"
	PL5hrFg  = "\033[38;2;25;50;55m"
	PL5hrTxt = "\033[38;2;100;200;210m"

	// 7day segment: dark purple
	PL7dayBg  = "\033[48;2;45;30;55m"
	PL7dayFg  = "\033[38;2;45;30;55m"
	PL7dayTxt = "\033[38;2;180;140;210m"
)

func (t *OnelinePowerlineTheme) Render(data StatusData) string {
	var sb strings.Builder

	arrow := "\ue0b0" // Powerline right arrow

	// Segment 1: Model + Version
	modelColor, _ := GetModelConfig(data.ModelType)
	sb.WriteString(PLModelBg)
	sb.WriteString(fmt.Sprintf(" %s%s%s", modelColor, data.ModelName, Reset))
	sb.WriteString(PLModelBg)
	sb.WriteString(fmt.Sprintf(" %s%s%s ", ColorDim, data.Version, Reset))
	sb.WriteString(PLPathBg)
	sb.WriteString(PLModelFg)
	sb.WriteString(arrow)
	sb.WriteString(Reset)

	// Segment 2: Path
	sb.WriteString(PLPathBg)
	sb.WriteString(fmt.Sprintf(" %s%s%s ", PLPathTxt, ShortenPath(data.ProjectPath, 20), Reset))

	// Segment 3: Git (or skip to 5hr)
	if data.GitBranch != "" {
		sb.WriteString(PLGitBg)
		sb.WriteString(PLPathFg)
		sb.WriteString(arrow)
		sb.WriteString(Reset)
		sb.WriteString(PLGitBg)
		sb.WriteString(fmt.Sprintf(" %s%s", PLGitTxt, data.GitBranch))
		if data.GitStaged > 0 {
			sb.WriteString(fmt.Sprintf(" +%d", data.GitStaged))
		}
		if data.GitDirty > 0 {
			sb.WriteString(fmt.Sprintf(" ~%d", data.GitDirty))
		}
		sb.WriteString(fmt.Sprintf(" %s", Reset))
		sb.WriteString(PL5hrBg)
		sb.WriteString(PLGitFg)
	} else {
		sb.WriteString(PL5hrBg)
		sb.WriteString(PLPathFg)
	}
	sb.WriteString(arrow)
	sb.WriteString(Reset)

	// Segment 4: 5hr
	color5 := t.pctColor(data.API5hrPercent)
	sb.WriteString(PL5hrBg)
	sb.WriteString(fmt.Sprintf(" %s5h%s %s%d%%%s", ColorDim, Reset+PL5hrBg, color5, data.API5hrPercent, Reset+PL5hrBg))
	if data.API5hrTimeLeft != "" {
		sb.WriteString(fmt.Sprintf(" %s%s%s", ColorDim, data.API5hrTimeLeft, Reset+PL5hrBg))
	}
	sb.WriteString(" ")
	sb.WriteString(PL7dayBg)
	sb.WriteString(PL5hrFg)
	sb.WriteString(arrow)
	sb.WriteString(Reset)

	// Segment 5: 7day
	color7 := t.pctColor(data.API7dayPercent)
	sb.WriteString(PL7dayBg)
	sb.WriteString(fmt.Sprintf(" %s7d%s %s%d%%%s", ColorDim, Reset+PL7dayBg, color7, data.API7dayPercent, Reset+PL7dayBg))
	if data.API7dayTimeLeft != "" {
		sb.WriteString(fmt.Sprintf(" %s%s%s", ColorDim, data.API7dayTimeLeft, Reset+PL7dayBg))
	}
	sb.WriteString(" ")
	sb.WriteString(Reset)
	sb.WriteString(PL7dayFg)
	sb.WriteString(arrow)
	sb.WriteString(Reset)

	sb.WriteString("\n")
	return sb.String()
}

func (t *OnelinePowerlineTheme) pctColor(pct int) string {
	if pct < 50 {
		return ColorBrightGreen
	} else if pct < 75 {
		return ColorBrightYellow
	}
	return ColorRed
}
