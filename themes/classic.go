package themes

import (
	"fmt"
	"strings"
)

// ClassicTheme original classic style
type ClassicTheme struct{}

func init() {
	RegisterTheme(&ClassicTheme{})
}

func (t *ClassicTheme) Name() string {
	return "classic"
}

func (t *ClassicTheme) Description() string {
	return "Original classic: maintains the original statusline layout style"
}

func (t *ClassicTheme) Render(data StatusData) string {
	var sb strings.Builder

	const tableWidth = 88
	const colWidth = 35

	// Line 1: Path + Git (left) + Model (right-aligned to table width)
	leftPart := t.formatPathGit(data)
	modelPart := t.formatModelShort(data)
	leftTargetWidth := tableWidth - 13
	sb.WriteString(Reset)
	sb.WriteString(PadRight(leftPart, leftTargetWidth))
	sb.WriteString(modelPart)
	sb.WriteString(Reset)
	sb.WriteString("\n")

	// Line 2: Cost
	sessCost := fmt.Sprintf("ses %s%s%s", ColorGreen, FormatCost(data.SessionCost), Reset)
	dayCost := fmt.Sprintf("day %s%s%s", ColorGold, FormatCost(data.DayCost), Reset)
	monCost := fmt.Sprintf("mon %s%s%s", ColorPurple, FormatCost(data.MonthCost), Reset)
	wkCost := fmt.Sprintf("week %s%s%s", ColorBlue, FormatCost(data.WeekCost), Reset)
	burnRate := fmt.Sprintf("avg %s%s/h%s", ColorRed, FormatCost(data.BurnRate), Reset)
	cacheStr := t.formatCachePercent(data)
	costCol1 := sessCost + "  " + dayCost + "  " + monCost
	costCol2 := wkCost + "  " + burnRate + "  " + cacheStr

	sb.WriteString(fmt.Sprintf("%s├─ %-9s │ %s│ %s│%s\n",
		ColorDim, "Cost", PadRight(costCol1, colWidth), PadRight(costCol2, colWidth), Reset))

	// Line 3: Stats + Context bar
	tokenStr := fmt.Sprintf("tok %s%s%s", ColorPurple, FormatTokensFixed(data.TokenCount, 6), Reset)
	msgStr := fmt.Sprintf("msg %s%4d%s", ColorCyan, data.MessageCount, Reset)
	timeStr := fmt.Sprintf("time %s", data.SessionTime)
	ctxBar := t.formatContextBar(data)
	statsCol1 := tokenStr + "  " + msgStr + "    " + timeStr
	statsCol2 := ctxBar

	sb.WriteString(fmt.Sprintf("%s├─ %-9s │ %s│ %s│%s\n",
		ColorDim, "Stats", PadRight(statsCol1, colWidth), PadRight(statsCol2, colWidth), Reset))

	// Line 4: API Limits
	api5hr := t.formatAPILimit(data.API5hrPercent, data.API5hrTimeLeft, "5hr")
	api7day := t.formatAPILimit(data.API7dayPercent, data.API7dayTimeLeft, "7day")

	sb.WriteString(fmt.Sprintf("%s└─ %-9s │ %s│ %s│%s\n",
		ColorDim, "API Limit", PadRight(api5hr, colWidth), PadRight(api7day, colWidth), Reset))

	return sb.String()
}

func (t *ClassicTheme) formatPathGit(data StatusData) string {
	path := fmt.Sprintf("📂 %s", data.ProjectPath)

	git := ""
	if data.GitBranch != "" {
		git = fmt.Sprintf("  %s⚡ %s%s", ColorCyan, data.GitBranch, Reset)
		if data.GitStaged > 0 || data.GitDirty > 0 || data.GitAhead > 0 || data.GitBehind > 0 || data.GitStash > 0 || data.GitSHA != "" {
			var status []string
			if data.GitStaged > 0 {
				status = append(status, fmt.Sprintf("%s+%d%s", ColorGreen, data.GitStaged, Reset))
			}
			if data.GitDirty > 0 {
				status = append(status, fmt.Sprintf("%s~%d%s", ColorOrange, data.GitDirty, Reset))
			}
			if ab := strings.TrimSpace(FormatGitExtras(data, ColorGreen, ColorOrange, Dim)); ab != "" {
				status = append(status, ab)
			}
			git += "  " + strings.Join(status, " ")
		}
	}

	return path + git
}

func (t *ClassicTheme) formatModelShort(data StatusData) string {
	modelColor, modelIcon := GetModelConfig(data.ModelType)
	return fmt.Sprintf("[%s%s %s%s]", modelColor, modelIcon, data.ModelName, Reset)
}

func (t *ClassicTheme) formatCachePercent(data StatusData) string {
	color := ColorGreen
	if data.CacheHitRate < 40 {
		color = ColorOrange
	} else if data.CacheHitRate < 70 {
		color = ColorYellow
	}
	return fmt.Sprintf("hit %s%3d%%%s", color, data.CacheHitRate, Reset)
}

func (t *ClassicTheme) formatContextBar(data StatusData) string {
	bar := GenerateBar(data.ContextPercent, 14, "█", "░", GetContextColor(data.ContextPercent), ColorGray)
	color := GetContextColor(data.ContextPercent)
	return fmt.Sprintf("Ctx  %s %s%3d%%%s %s", bar, color, data.ContextPercent, Reset, FormatNumber(data.ContextUsed))
}

func (t *ClassicTheme) formatAPILimit(percent int, timeLeft, label string) string {
	bar := GenerateBar(percent, 14, "█", "░", getAPIColor(percent), ColorGray)
	color := getAPIColor(percent)
	return fmt.Sprintf("%s %s %s%3d%%%s (%s)", label, bar, color, percent, Reset, timeLeft)
}

func getAPIColor(percent int) string {
	if percent < 50 {
		return ColorGreen
	} else if percent < 75 {
		return ColorYellow
	} else if percent < 90 {
		return ColorOrange
	}
	return ColorRed
}
