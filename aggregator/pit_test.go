package aggregator

import (
	"testing"

	"github.com/nstandage/f1-go-cli-app/model"
)

func makeStoreWithDrivers() *Store {
	return &Store{
		Drivers: map[uint]*Driver{
			1: {Number: 1, Info: &model.Driver{DriverNumber: 1, NameAcronym: "VER"}},
			2: {Number: 2, Info: &model.Driver{DriverNumber: 2, NameAcronym: "HAM"}},
		},
	}
}

func TestUpdatePit_AddsEntryToFront(t *testing.T) {
	s := makeStoreWithDrivers()
	s.updatePit(&model.Pit{DriverNumber: 1, StopDuration: 2.5})
	if len(s.RecentPits) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(s.RecentPits))
	}
	if s.RecentPits[0].DriverAcronym != "VER" {
		t.Errorf("expected VER, got %s", s.RecentPits[0].DriverAcronym)
	}
	if s.RecentPits[0].StopDuration != 2.5 {
		t.Errorf("expected 2.5, got %f", s.RecentPits[0].StopDuration)
	}
}

func TestUpdatePit_NewestFirst(t *testing.T) {
	s := makeStoreWithDrivers()
	s.updatePit(&model.Pit{DriverNumber: 1, StopDuration: 2.5})
	s.updatePit(&model.Pit{DriverNumber: 2, StopDuration: 3.1})
	if s.RecentPits[0].DriverAcronym != "HAM" {
		t.Errorf("expected HAM at front, got %s", s.RecentPits[0].DriverAcronym)
	}
	if s.RecentPits[1].DriverAcronym != "VER" {
		t.Errorf("expected VER at index 1, got %s", s.RecentPits[1].DriverAcronym)
	}
}

func TestUpdatePit_CapsAtEight(t *testing.T) {
	s := makeStoreWithDrivers()
	for i := 0; i < 10; i++ {
		s.updatePit(&model.Pit{DriverNumber: 1, StopDuration: float64(i)})
	}
	if len(s.RecentPits) != 8 {
		t.Errorf("expected 8 entries, got %d", len(s.RecentPits))
	}
	// newest stop was StopDuration=9 — must be at front
	if s.RecentPits[0].StopDuration != 9.0 {
		t.Errorf("expected 9.0 at front, got %f", s.RecentPits[0].StopDuration)
	}
}

func TestUpdatePit_UnknownDriverSkipped(t *testing.T) {
	s := makeStoreWithDrivers()
	s.updatePit(&model.Pit{DriverNumber: 99, StopDuration: 2.5})
	if len(s.RecentPits) != 0 {
		t.Errorf("expected 0 entries for unknown driver, got %d", len(s.RecentPits))
	}
}
