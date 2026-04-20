package view

import (
	"fmt"
	"image/color"
)

func PitStops(stops []float64) string {
	str := ""

	for _, stop := range stops {
		c := getPitColor(stop)
		s := fmt.Sprintf("%.2f", stop)
		str = str + defaultTextStyle(s, c) + "\n"
	}

	return defaultBorderStyle().Width(22).Height(14).Render(str)
}

func getPitColor(stop float64) color.Color {
	switch {
	case stop < 3.1:
		return pitStopFastColor
	case stop < 3.5:
		return pitStopAverageColor
	default:
		return pitStopSlowColor

	}
}
