package aggregator

import "github.com/nstandage/f1-go-cli-app/model"

// getActiveStint returns the stint with the highest LapStart that is <= currentLap.
// Returns nil if no stint has started yet.
func getActiveStint(stints []model.Stint, currentLap uint) *model.Stint {
	var active *model.Stint
	for i := range stints {
		s := &stints[i]
		if s.LapStart <= currentLap {
			if active == nil || s.LapStart > active.LapStart {
				active = s
			}
		}
	}
	return active
}

// tireAge returns the total age of a tyre set at the given lap,
// including any age the tyres had before the stint began.
func tireAge(stint *model.Stint, currentLap uint) uint {
	return stint.TyreAgeAtStart + (currentLap - stint.LapStart)
}
