package view

import (
	"fmt"
	"image/color"
	"charm.land/lipgloss/v2"
	"github.com/nstandage/f1-go-cli-app/model"
)

func PositionsColumn() string {
	return lipgloss.NewStyle().
		Padding(0, 0, 0, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Right,
			getPositionNumbers(22)...,
		),
		)
}

func getPositionNumbers(n int) []string {
	rows := make([]string, n)
	for i := range n {
		string := fmt.Sprintf("%2d", i+1)
		rows[i] = defaultTextStyle(string, title2Color)
	}
	return rows
}

func DefaultColumn(i []string) string {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 6).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styleStrings(i, title2Color)...,
			),
		)
}

func DriverColumn(drivers []model.DriverSnapshot) string {
	styledDrivers := make([]string, len(drivers))
	for i, d := range drivers {
		c := title2Color
		if d.IsFastestLap {
			c = bestOverallSectorColor
		}
		styledDrivers[i] = defaultTextStyle(d.Name, c)
	}

	return lipgloss.NewStyle().
		Margin(0, 0, 0, 6).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styledDrivers...
			),
		)
}


func LastLapColumn(laps []string, isPitOut []bool) string {
	styled := make([]string, len(laps))
	for i, lap := range laps {
		var c color.Color = title2Color
		if i < len(isPitOut) && isPitOut[i] {
			c = pitLaneSectorColor
		}
		styled[i] = defaultTextStyle(lap, c)
	}
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 3).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, styled...),
		)
}

func PitColumn(i []string) string {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 10).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styleStrings(i, title2Color)...,
			),
		)
}

func TireColumn(tires []string) string {
	styledTires := make([]string, len(tires))

	for i, tire := range tires {
		c := getTireColor(tire)
		styledTires[i] = defaultTextStyle(tire, c)
	}

	return lipgloss.NewStyle().
		Margin(0, 0, 0, 5).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				styledTires...,
			),
		)
}

func getTireColor(s string) color.Color {
	switch s {
	case "SOFT":
		return softTireColor
	case "MEDIUM":
		return mediumTireColor
	case "HARD":
		return hardTireColor
	case "INT":
		return intTireColor
	case "WET":
		return wetTireColor
	default:
		return mediumTireColor
	}
}

func TireAgeColumn(i []string) string {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 4).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Right,
				styleStrings(i, title2Color)...,
			),
		)
}
