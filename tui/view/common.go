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

	pitStopFastColor    = lipgloss.Color("#00e676")
	pitStopAverageColor = lipgloss.Color("#ffd600")
	pitStopSlowColor    = lipgloss.Color("#e3504d")
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
