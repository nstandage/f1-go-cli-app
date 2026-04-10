package view

import "charm.land/lipgloss/v2"

func LegendBar() string {
	return defaultBorderStyle().Render(
		defaultTextStyle("Sector Colors:", lipgloss.White),
		defaultTextStyle(fullShadeBlock, bestOverallSectorColor),
		defaultTextStyle("Overall Best", titleDarkColor),
		defaultTextStyle(fullShadeBlock, bestPersonalSectorColor),
		defaultTextStyle("Personal Best", titleDarkColor),
		defaultTextStyle(fullShadeBlock, slowSectorColor),
		defaultTextStyle("Slower", titleDarkColor),
		defaultDivider,
		defaultTextStyle("Pit Stop Colors:", lipgloss.White),
		defaultTextStyle(fullShadeBlock, pitStopFastColor),
		defaultTextStyle("<3.1s", titleDarkColor),
		defaultTextStyle(fullShadeBlock, pitStopAverageColor),
		defaultTextStyle("<3.5", titleDarkColor),
		defaultTextStyle(fullShadeBlock, pitStopSlowColor),
		defaultTextStyle(">3.5", titleDarkColor),
		defaultDivider,
		defaultTextStyle("q quit", titleDarkColor),
	)
}
