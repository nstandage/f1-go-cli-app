package view

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

func Topbar(sectors []int) string {
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
		sectorTitle("S1", sectors[0]),
		sectorTitle("S2", sectors[1]),
		sectorTitle("S3", sectors[2]),
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

func sectorTitle(sector string, miniSectors int) string {
	// This is ugly. minisectors - number of characters that are in the first Sprintf. times 2 because that's how many characters per mini sector I'm showing. (lap.go)
	title := fmt.Sprintf("--%v", sector) + strings.Repeat("-", (miniSectors-2)*2) //
	// title := "----------"
	return lipgloss.
		NewStyle().
		Margin(0, 1).
		Render(defaultTextStyle(title, titleDarkColor))
}
