package view

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

func Laps(drivers [][][]uint) string {
	rows := make([]string, len(drivers))
	for i, sectors := range drivers {
		row := ""
		for _, sector := range sectors {
			for _, seg := range sector {
				row += defaultTextStyle(lowerHalfBlock, miniSectorColor(seg)) + " "
			}
			row += "   "
		}
		rows[i] = row
	}
	return lipgloss.NewStyle().Margin(0, 0, 0, 7).Render(
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func miniSectorColor(i uint) color.Color {
	switch i {
	case 2048:
		return slowSectorColor
	case 2049:
		return bestPersonalSectorColor
	case 2051:
		return bestOverallSectorColor
	case 2064:
		return pitLaneSectorColor
	default:
		return futureSectorColor
	}
}
