package view

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
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

func LastLapColumn(i []string) string {
	return lipgloss.NewStyle().
		Margin(0, 0, 0, 3).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				styleStrings(i, title2Color)...,
			),
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
