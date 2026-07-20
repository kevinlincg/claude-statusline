package themes

import (
	"fmt"
	"strings"
)

// MatrixTheme matrix hacker style
type MatrixTheme struct{}

func init() {
	RegisterTheme(&MatrixTheme{})
}

func (t *MatrixTheme) Name() string {
	return "matrix"
}

func (t *MatrixTheme) Description() string {
	return "Matrix hacker: green terminal style"
}

const (
	MatrixGreen     = "\033[38;2;0;255;0m"
	MatrixDarkGreen = "\033[38;2;0;180;0m"
	MatrixBg        = "\033[48;2;0;30;0m"
)

func (t *MatrixTheme) Render(data StatusData) string {
	var sb strings.Builder

	const width = 95

	// Top border
	sb.WriteString(MatrixGreen)
	sb.WriteString("░▒▓")
	sb.WriteString(strings.Repeat("█", width-4))
	sb.WriteString("▓▒░")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Line 1: Model + Version | Path + Git
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	update := ""
	if data.UpdateAvailable {
		update = fmt.Sprintf(" %s⬆%s", ColorNeonOrange, Reset)
	}

	line1 := fmt.Sprintf(" %s$>%s %s%s%s%s%s %s%s%s%s  %s│%s  %s📂 %s%s",
		MatrixDarkGreen, Reset,
		modelColor, Bold, modelIcon, data.ModelName, Reset,
		ColorNeonGreen, data.Version, Reset, update,
		MatrixDarkGreen, Reset,
		ColorYellow, data.ProjectPath, Reset)
	if data.GitBranch != "" {
		line1 += fmt.Sprintf("  %s⚡%s%s", MatrixGreen, data.GitBranch, Reset)
		if data.GitStaged > 0 {
			line1 += fmt.Sprintf(" %s+%d%s", ColorGreen, data.GitStaged, Reset)
		}
		if data.GitDirty > 0 {
			line1 += fmt.Sprintf(" %s~%d%s", ColorOrange, data.GitDirty, Reset)
		}
		line1 += FormatGitExtras(data, ColorGreen, ColorOrange, Dim)
	}

	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line1, width-2))
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Line 2: Session Stats | Cost
	line2 := fmt.Sprintf(" %s$>%s %s%5s%s tok  %s%3d%s msg  %s%6s%s  %s│%s  %s%s%s ses  %s%s%s day  %s%s/h%s  %s%d%%hit%s",
		MatrixDarkGreen, Reset,
		ColorPurple, FormatTokens(data.TokenCount), Reset,
		ColorCyan, data.MessageCount, Reset,
		ColorSilver, data.SessionTime, Reset,
		MatrixDarkGreen, Reset,
		MatrixGreen, FormatCost(data.SessionCost), Reset,
		ColorYellow, FormatCost(data.DayCost), Reset,
		ColorRed, FormatCost(data.BurnRate), Reset,
		ColorGreen, data.CacheHitRate, Reset)

	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line2, width-2))
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Line 3: Progress bars
	color1, _ := GetBarColor(data.ContextPercent)
	color5, _ := GetBarColor(data.API5hrPercent)
	color7, _ := GetBarColor(data.API7dayPercent)

	line3 := fmt.Sprintf(" %s$>%s %sCtx%s %s %s%3d%%%s  %s│%s  %s5hr%s %s %s%3d%%%s %s%s%s  %s│%s  %s7dy%s %s %s%3d%%%s %s%s%s",
		MatrixDarkGreen, Reset,
		ColorLabelDim, Reset,
		GenerateGlowBar(data.ContextPercent, 14, color1, MatrixBg),
		color1, data.ContextPercent, Reset,
		MatrixDarkGreen, Reset,
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API5hrPercent, 10, color5, MatrixBg),
		color5, data.API5hrPercent, Reset,
		ColorDim, data.API5hrTimeLeft, Reset,
		MatrixDarkGreen, Reset,
		ColorLabelDim, Reset,
		GenerateGlowBar(data.API7dayPercent, 10, color7, MatrixBg),
		color7, data.API7dayPercent, Reset,
		ColorDim, data.API7dayTimeLeft, Reset)

	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString(PadRight(line3, width-2))
	sb.WriteString(MatrixGreen)
	sb.WriteString("▓")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Bottom border
	sb.WriteString(MatrixGreen)
	sb.WriteString("░▒▓")
	sb.WriteString(strings.Repeat("█", width-4))
	sb.WriteString("▓▒░")
	sb.WriteString(Reset)
	sb.WriteString("\n")

	return sb.String()
}
