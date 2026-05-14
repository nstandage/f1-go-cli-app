package aggregator

import (
	"testing"

	"github.com/nstandage/f1-go-cli-app/model"
)

func TestGetActiveStint_NoStints(t *testing.T) {
	result := getActiveStint(nil, 5)
	if result != nil {
		t.Errorf("expected nil for empty stints, got %v", result)
	}
}

func TestGetActiveStint_CurrentLapZero(t *testing.T) {
	stints := []model.Stint{{LapStart: 1, LapEnd: 20, Compound: "SOFT"}}
	result := getActiveStint(stints, 0)
	if result != nil {
		t.Errorf("expected nil for lap 0, got %v", result)
	}
}

func TestGetActiveStint_SingleStint(t *testing.T) {
	stints := []model.Stint{{LapStart: 1, LapEnd: 20, Compound: "SOFT"}}
	result := getActiveStint(stints, 10)
	if result == nil {
		t.Fatal("expected stint, got nil")
	}
	if result.Compound != "SOFT" {
		t.Errorf("expected SOFT, got %v", result.Compound)
	}
}

func TestGetActiveStint_ReturnsLatestActiveStint(t *testing.T) {
	stints := []model.Stint{
		{LapStart: 1, LapEnd: 20, Compound: "SOFT", StintNumber: 1},
		{LapStart: 21, LapEnd: 45, Compound: "MEDIUM", StintNumber: 2},
	}
	result := getActiveStint(stints, 25)
	if result == nil {
		t.Fatal("expected stint, got nil")
	}
	if result.Compound != "MEDIUM" {
		t.Errorf("expected MEDIUM, got %v", result.Compound)
	}
}

func TestGetActiveStint_BeforeSecondStintStarts(t *testing.T) {
	stints := []model.Stint{
		{LapStart: 1, LapEnd: 20, Compound: "SOFT", StintNumber: 1},
		{LapStart: 21, LapEnd: 45, Compound: "MEDIUM", StintNumber: 2},
	}
	result := getActiveStint(stints, 15)
	if result == nil {
		t.Fatal("expected stint, got nil")
	}
	if result.Compound != "SOFT" {
		t.Errorf("expected SOFT, got %v", result.Compound)
	}
}

func TestTireAge_NoPreUsedTyres(t *testing.T) {
	stint := &model.Stint{LapStart: 1, TyreAgeAtStart: 0}
	age := tireAge(stint, 10)
	if age != 9 {
		t.Errorf("expected 9, got %v", age)
	}
}

func TestTireAge_WithPreUsedTyres(t *testing.T) {
	stint := &model.Stint{LapStart: 15, TyreAgeAtStart: 5}
	age := tireAge(stint, 20)
	if age != 10 {
		t.Errorf("expected 10, got %v", age)
	}
}

func TestTireAge_AtLapStart(t *testing.T) {
	stint := &model.Stint{LapStart: 10, TyreAgeAtStart: 3}
	age := tireAge(stint, 10)
	if age != 3 {
		t.Errorf("expected 3 (TyreAgeAtStart only), got %v", age)
	}
}

func TestGetActiveStint_UnorderedStints(t *testing.T) {
	// stints provided in reverse order — implementation must still find correct active stint
	stints := []model.Stint{
		{LapStart: 21, LapEnd: 45, Compound: "MEDIUM", StintNumber: 2},
		{LapStart: 1, LapEnd: 20, Compound: "SOFT", StintNumber: 1},
	}
	result := getActiveStint(stints, 15)
	if result == nil {
		t.Fatal("expected stint, got nil")
	}
	if result.Compound != "SOFT" {
		t.Errorf("expected SOFT (first stint), got %v", result.Compound)
	}
}
