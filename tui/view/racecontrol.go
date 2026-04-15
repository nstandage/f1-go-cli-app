package view

func RaceControl(msgs []string) string {
	str := ""

	for _, msg := range msgs {
		str = str + raceControlBullet + msg + "\n"
	}

	return defaultBorderStyle().Width(22).Height(18).Render(str)
}
