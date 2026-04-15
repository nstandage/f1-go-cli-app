package view

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

func Laps(sectors [][]int) string {
	str := ""

	for _, sector := range sectors {
		for _, miniSector := range sector {
			color := miniSectorColor(miniSector)
			block := strings.Repeat(fullShadeBlock, 2)
			str = str + defaultTextStyle(block, color)
		}
		str = str + "   "

		//Number of spaces filling the margin gap in sectorTitle()
		// secString = secString + strings.Repeat(lightShadeBlock+lightShadeBlock, s) + "   "
	}

	return lipgloss.NewStyle().Margin(0, 0, 0, 7).Render(
		str,
	)
}

func miniSectorColor(i int) color.Color {
	switch i {
	case 2048: // Yellow sector
		return slowSectorColor
	case 2049: // Green sector
		return bestPersonalSectorColor
	case 2051: // Purple sector
		return bestOverallSectorColor
	case 2064: // Pitlane
		return pitLaneSectorColor
	default:
		return futureSectorColor
	}
}
