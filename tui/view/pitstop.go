package view

import (
	"fmt"
	"image/color"

	"github.com/nstandage/f1-go-cli-app/model"
)

func PitStops(stops []model.PitStopEntry) string {
	str := ""
	for _, stop := range stops {
		c := getPitColor(stop.StopDuration)
		row := fmt.Sprintf("%-4s %5.2f", stop.DriverAcronym, stop.StopDuration)
		str = str + defaultTextStyle(row, c) + "\n"
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
