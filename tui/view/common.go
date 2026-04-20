package view

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Colors
var (
	bestOverallSectorColor  = lipgloss.Color("#ce93d8")
	bestPersonalSectorColor = lipgloss.Color("#00e676")
	slowSectorColor         = lipgloss.Color("#ffd600")
	futureSectorColor       = lipgloss.Color("#3C3C3C")
	pitLaneSectorColor      = lipgloss.Color("#e3504d")

	pitStopFastColor    = lipgloss.Color("#00e676")
	pitStopAverageColor = lipgloss.Color("#ffd600")
	pitStopSlowColor    = lipgloss.Color("#e3504d")

	softTireColor   = lipgloss.Color("#e3504d")
	mediumTireColor = lipgloss.Color("#ffd600")
	hardTireColor   = lipgloss.BrightWhite
	intTireColor    = lipgloss.Color("#00e676")
	wetTireColor    = lipgloss.Color("#29CFE6")
)

var (
	title1Color    = lipgloss.Color("#29CFE6")
	title2Color    = lipgloss.BrightWhite
	title3Color    = lipgloss.Color("#C78ED0")
	titleDarkColor = lipgloss.Color("#666666")
	borderColor    = lipgloss.Color("#3C3C3C")
)

// Text
var defaultDivider = defaultTextStyle(" • ", titleDarkColor)
var raceControlBullet = defaultTextStyle("• ", slowSectorColor)

const (
	lightShadeBlock  = "░" // U+2591
	mediumShadeBlock = "▒" // U+2592
	darkShadeBlock   = "▓" // U+2593
	fullShadeBlock   = "█" // U+2588
)

// Styles
func defaultBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor)
}

func defaultTextStyle(s string, c color.Color) string {
	var style = lipgloss.NewStyle().
		Foreground(c).
		Bold(true)

	return style.Render(s)
}

func styleStrings(rows []string, c color.Color) []string {
	styledRows := make([]string, len(rows))
	for i, row := range rows {
		styledRows[i] = defaultTextStyle(row, c)
	}
	return styledRows
}

func Spacer(size int) string {
	return lipgloss.NewStyle().Height(size).Render("")
}
