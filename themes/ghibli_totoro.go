package themes

import (
	"fmt"
	"strings"
)

// TotoroTheme My Neighbor Totoro forest spirit style
type TotoroTheme struct{}

func init() {
	RegisterTheme(&TotoroTheme{})
}

func (t *TotoroTheme) Name() string {
	return "totoro"
}

func (t *TotoroTheme) Description() string {
	return "Totoro: My Neighbor Totoro forest spirit style"
}

const (
	TotoroGreen     = "\033[38;2;144;238;144m"
	TotoroDarkGreen = "\033[38;2;34;139;34m"
	TotoroBrown     = "\033[38;2;139;90;43m"
	TotoroSky       = "\033[38;2;135;206;235m"
	TotoroCream     = "\033[38;2;255;253;208m"
	TotoroGray      = "\033[38;2;128;128;128m"
	TotoroYellow    = "\033[38;2;255;255;150m"
	TotoroWhite     = "\033[38;2;255;255;255m"
)

func (t *TotoroTheme) Render(data StatusData) string {
	var sb strings.Builder

	// Forest canopy top
	sb.WriteString(TotoroDarkGreen + "  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + Reset + "\n")

	// Peaceful header
	modelColor, modelIcon := GetModelConfig(data.ModelType)

	title := fmt.Sprintf("    %s🌳%s %s%s%s %s %s%s%s",
		TotoroGreen, Reset,
		modelColor, modelIcon, data.ModelName,
		TotoroGray, data.Version, Reset,
		func() string {
			if data.UpdateAvailable {
				return fmt.Sprintf(" %s✨new%s", TotoroYellow, Reset)
			}
			return ""
		}())

	sb.WriteString(title + "\n")

	sb.WriteString(TotoroDarkGreen + "  ─────────────────────────────────────────────────────────────────────────────────────────" + Reset + "\n")

	// Forest path (project)
	gitInfo := ""
	if data.GitBranch != "" {
		gitInfo = fmt.Sprintf("  %s🍃%s%s", TotoroGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			gitInfo += fmt.Sprintf(" %s+%d%s", TotoroGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			gitInfo += fmt.Sprintf(" %s~%d%s", TotoroYellow, data.GitDirty, Reset)
		}
		gitInfo += FormatGitExtras(data, TotoroGreen, TotoroYellow, Dim)
	}

	line1 := fmt.Sprintf("    %s🏠 Path:%s %s%s",
		TotoroBrown, Reset, ShortenPath(data.ProjectPath, 40), gitInfo)
	sb.WriteString(line1 + "\n")

	sb.WriteString("\n")

	// Gentle progress indicators
	spiritBar := t.generateTotoroBar(data.ContextPercent, 20)
	line2 := fmt.Sprintf("    %s🌱 Spirit Energy%s  %s  %s%3d%%%s",
		TotoroGreen, Reset, spiritBar, TotoroGreen, data.ContextPercent, Reset)
	sb.WriteString(line2 + "\n")

	catbusBar := t.generateTotoroBar(100-data.API5hrPercent, 20)
	line3 := fmt.Sprintf("    %s🐱 Catbus Fuel%s    %s  %s%3d%%%s  %s%s%s",
		TotoroYellow, Reset, catbusBar, TotoroYellow, 100-data.API5hrPercent, Reset,
		TotoroGray, data.API5hrTimeLeft, Reset)
	sb.WriteString(line3 + "\n")

	acornBar := t.generateTotoroBar(100-data.API7dayPercent, 20)
	line4 := fmt.Sprintf("    %s🌰 Acorn Storage%s  %s  %s%3d%%%s  %s%s%s",
		TotoroBrown, Reset, acornBar, TotoroBrown, 100-data.API7dayPercent, Reset,
		TotoroGray, data.API7dayTimeLeft, Reset)
	sb.WriteString(line4 + "\n")

	sb.WriteString("\n")

	// Stats in soft style
	line5 := fmt.Sprintf("    %s☁️%s %s tok  %s🕐%s %s  %s💫%s %d msg  %s🍂%s $%s  %s☀️%s $%s/day  %s🌸%s %d%% luck",
		TotoroSky, Reset, FormatTokens(data.TokenCount),
		TotoroGray, Reset, data.SessionTime,
		TotoroYellow, Reset, data.MessageCount,
		TotoroCream, Reset, FormatCost(data.SessionCost),
		TotoroYellow, Reset, FormatCost(data.DayCost),
		TotoroGreen, Reset, data.CacheHitRate)
	sb.WriteString(line5 + "\n")

	// Forest floor
	sb.WriteString(TotoroDarkGreen + "  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + Reset + "\n")

	return sb.String()
}

func (t *TotoroTheme) generateTotoroBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * width / 100
	empty := width - filled

	// Soft nature colors based on level
	fillChar := "●"
	emptyChar := "○"
	color := TotoroGreen
	if percent > 75 {
		color = TotoroYellow
	}
	if percent > 90 {
		color = TotoroBrown
	}

	var bar strings.Builder
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat(fillChar, filled))
		bar.WriteString(Reset)
	}
	if empty > 0 {
		bar.WriteString(TotoroGray)
		bar.WriteString(strings.Repeat(emptyChar, empty))
		bar.WriteString(Reset)
	}
	return bar.String()
}
