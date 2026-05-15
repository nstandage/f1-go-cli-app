# Live Pit Stop Feed Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the static pit stop list in the TUI with a live feed showing the last 8 pit stops (newest first), each displaying driver acronym and stop duration.

**Architecture:** Add a `PitStopEntry` struct to the snapshot model, maintain a capped prepend-list of recent pits in the aggregator store, and thread it through the snapshot into the view. The existing `Pit` event pipeline (fetch → replay → engine.handle → store.updatePit) already works; only the store mutation and snapshot output need changing.

**Tech Stack:** Go, Bubbletea v2, Lipgloss v2

---

## File Map

| File | Change |
|------|--------|
| `model/snapshot.go` | Add `PitStopEntry` struct; add `RecentPitStops []PitStopEntry` to `Snapshot` |
| `aggregator/store.go` | Add `RecentPits []PitStopEntry` field; update `updatePit` to prepend + cap at 8 |
| `aggregator/pit_test.go` | New file — tests for `updatePit` behaviour |
| `aggregator/engine.go` | Populate `RecentPitStops` in `GetSnapshot` |
| `tui/view/pitstop.go` | Change signature to `[]model.PitStopEntry`; update render logic |
| `tui/model.go` | Remove hardcoded slice; pass `snapshot.RecentPitStops` to `view.PitStops` |

---

### Task 1: Add `PitStopEntry` to the snapshot model

**Files:**
- Modify: `model/snapshot.go`

- [ ] **Step 1: Add the struct and field**

In `model/snapshot.go`, add `PitStopEntry` and the new field on `Snapshot`. The full updated file:

```go
package model

import (
	"time"
)

type Snapshot struct {
	SessionBar      *SessionBarSnapShot
	RaceControlMsgs []string
	Drivers         []DriverSnapshot
	LastLap         []string
	LastLapIsPitOut []bool
	Intervals       []string
	GapsToLeaders   []string
	PitCounts       []string
	TireCompounds   []string
	TireAges        []string
	RecentPitStops  []PitStopEntry
}

type PitStopEntry struct {
	DriverAcronym string
	StopDuration  float64
}

type SessionBarSnapShot struct {
	EventName  string
	EventType  string
	FastestLap *FastestLapSnapshot
	CurrentLap uint
	TotalLaps  uint
	IsReplay   bool
	EventDate  time.Time
}

type FastestLapSnapshot struct {
	LapTime   string
	Driver    string
	LapNumber string
}

type DriverSnapshot struct {
	Name         string
	IsFastestLap bool
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./...
```

Expected: no output (clean build).

- [ ] **Step 3: Commit**

```bash
git add model/snapshot.go
git commit -m "Add PitStopEntry and RecentPitStops to Snapshot model"
```

---

### Task 2: Update the store to track recent pit stops

**Files:**
- Modify: `aggregator/store.go`
- Create: `aggregator/pit_test.go`

- [ ] **Step 1: Write failing tests**

Create `aggregator/pit_test.go`:

```go
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./aggregator/... -run TestUpdatePit -v
```

Expected: compilation error — `s.RecentPits` undefined.

- [ ] **Step 3: Add `RecentPits` field to `Store`**

In `aggregator/store.go`, add the field to the `Store` struct:

```go
type Store struct {
	history      []model.Snapshot
	Drivers      map[uint]*Driver
	RaceControl  []model.RaceControl
	TotalLaps    uint
	CurrentLap   uint
	IsReplay     bool
	Session      *model.Session
	Meeting      *model.Meeting
	StartingGrid []model.StartingGrid
	FastestLap   *FastestLap
	Stints       map[uint][]model.Stint
	RecentPits   []model.PitStopEntry
}
```

- [ ] **Step 4: Update `updatePit` to prepend and cap**

Replace the existing `updatePit` method in `aggregator/store.go`:

```go
func (s *Store) updatePit(p *model.Pit) {
	driver, ok := s.Drivers[p.DriverNumber]
	if !ok {
		return
	}
	driver.PitCount++
	entry := model.PitStopEntry{
		DriverAcronym: driver.Info.NameAcronym,
		StopDuration:  p.StopDuration,
	}
	s.RecentPits = append([]model.PitStopEntry{entry}, s.RecentPits...)
	if len(s.RecentPits) > 8 {
		s.RecentPits = s.RecentPits[:8]
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test ./aggregator/... -run TestUpdatePit -v
```

Expected: all four `TestUpdatePit_*` tests PASS.

- [ ] **Step 6: Run all tests**

```bash
go test ./...
```

Expected: all tests pass, no failures.

- [ ] **Step 7: Commit**

```bash
git add aggregator/store.go aggregator/pit_test.go
git commit -m "Track last 8 pit stops in store, newest first"
```

---

### Task 3: Expose recent pit stops in the snapshot

**Files:**
- Modify: `aggregator/engine.go`

- [ ] **Step 1: Populate `RecentPitStops` in `GetSnapshot`**

In `aggregator/engine.go`, update the `snapshot` literal inside `GetSnapshot` to add the new field. The updated snapshot construction (lines 130–142):

```go
snapshot := model.Snapshot{
    SessionBar:      &sessionBar,
    RaceControlMsgs: e.getRaceControlMessages(),
    Drivers:         e.getDriverNames(),
    LastLap:         lastLap,
    LastLapIsPitOut: lastLapIsPitOut,
    Intervals:       e.getIntervals(),
    GapsToLeaders:   e.getGapToLeader(),
    PitCounts:       e.getPitCounts(),
    TireCompounds:   e.getTireCompounds(),
    TireAges:        e.getTireAges(),
    RecentPitStops:  e.store.RecentPits,
}
```

- [ ] **Step 2: Build and run all tests**

```bash
go build ./... && go test ./...
```

Expected: clean build, all tests pass.

- [ ] **Step 3: Commit**

```bash
git add aggregator/engine.go
git commit -m "Expose RecentPitStops in engine snapshot"
```

---

### Task 4: Update the pit stop view

**Files:**
- Modify: `tui/view/pitstop.go`

- [ ] **Step 1: Rewrite the view function**

Replace the entire contents of `tui/view/pitstop.go`:

```go
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
```

- [ ] **Step 2: Build**

```bash
go build ./...
```

Expected: compilation error — `tui/model.go` still passes `[]float64` to `PitStops`. That's expected; fix in next task.

---

### Task 5: Wire live data into the TUI

**Files:**
- Modify: `tui/model.go`

- [ ] **Step 1: Replace hardcoded pit stop slice**

In `tui/model.go`, remove the hardcoded slice and update the `view.PitStops` call. Find and replace this block (around lines 89–102):

Before:
```go
pitStops := []float64{
    3.0, 3.2, 3.8, 2.99, 3.12,
}

driverColumn := view.DriverColumn(drivers)
intervalColumn := view.DefaultColumn(intervals)
gapToLeaderColumn := view.DefaultColumn(gapToLeaders)
lastLapColumn := view.LastLapColumn(lastLap, snapshot.LastLapIsPitOut)
pitColumn := view.PitColumn(pits)
tiresColumn := view.TireColumn(tires)
tireAgeColumn := view.TireAgeColumn(tireAge)
laps := view.Laps(lapSectors)
raceControl := view.RaceControl(raceControlMessages)
pitStopView := view.PitStops(pitStops)
```

After:
```go
driverColumn := view.DriverColumn(drivers)
intervalColumn := view.DefaultColumn(intervals)
gapToLeaderColumn := view.DefaultColumn(gapToLeaders)
lastLapColumn := view.LastLapColumn(lastLap, snapshot.LastLapIsPitOut)
pitColumn := view.PitColumn(pits)
tiresColumn := view.TireColumn(tires)
tireAgeColumn := view.TireAgeColumn(tireAge)
laps := view.Laps(lapSectors)
raceControl := view.RaceControl(raceControlMessages)
pitStopView := view.PitStops(snapshot.RecentPitStops)
```

- [ ] **Step 2: Build and run all tests**

```bash
go build ./... && go test ./...
```

Expected: clean build, all tests pass.

- [ ] **Step 3: Commit**

```bash
git add tui/view/pitstop.go tui/model.go
git commit -m "Wire live pit stop feed into TUI"
```

---

### Task 6: Manual smoke test

- [ ] **Step 1: Run the app**

```bash
go run .
```

Watch the pit stop panel (bottom-right). During replay, each time a driver pits the panel should update: the new stop appears at the top with the driver's 3-letter acronym and stop time coloured green/yellow/red. Older stops shift down; the list never exceeds 8 rows.

- [ ] **Step 2: Verify empty state**

At the very start of the replay (before any pit stops), the panel should render with an empty body and just the border — no crash, no placeholder data.
