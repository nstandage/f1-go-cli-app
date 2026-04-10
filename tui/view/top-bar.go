package view

import (
	"charm.land/lipgloss/v2"
)

func Topbar() string {
	return defaultBorderStyle().Render(
		topTitle("Pos", 1),
		topTitle("Driver", 1),
		topTitle("Interval", 1),
		topTitle("To Lead", 2),
		topTitle("Last Lap", 2),
		defaultDivider,
		topTitle("Pits", 1),
		topTitle("Tire", 2),
		topTitle("Age", 1),
		defaultDivider,
		topTitle("--S1------", 1),
		topTitle("--S2------", 1),
		topTitle("--S3------", 1),
		defaultDivider,
		topTitle("Race Control", 1),
	)
}

func topTitle(title string, margin int) string {
	return lipgloss.
		NewStyle().
		Margin(0, margin).
		Render(defaultTextStyle(title, titleDarkColor))
}
