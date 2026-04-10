package view

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

func PositionsColumn() string {
	return lipgloss.NewStyle().
		Padding(0,0,0,2).
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
