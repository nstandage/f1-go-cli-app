# Pit Count, Tire Compound & Tire Age Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace hardcoded pit/tire placeholder data in the TUI with live values derived from the OpenF1 `/pit` and `/stints` API endpoints.

**Architecture:** Stints are fetched once as reference data and stored in the aggregator keyed by driver number; tire compound and age are computed per-snapshot from `CurrentLap`. Pit events are replayed through the existing event stream so `Driver.PitCount` increments at the correct race moment, and the pit data infrastructure is reusable for a future stop-duration feature.

**Tech Stack:** Go 1.26, `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, OpenF1 REST API.

---

## File Map

| File | Change |
|------|--------|
| `model/snapshot.go` | Add `PitCounts []int`, `TireCompounds []string`, `TireAges []string` to `Snapshot` |
| `aggregator/store.go` | Add `Stints map[uint][]model.Stint` to `Store`; add `PitCount uint` to `Driver` |
| `aggregator/tire.go` | **New file.** Package-level helpers: `getActiveStint`, `tireAge` |
| `aggregator/tire_test.go` | **New file.** Tests for `getActiveStint` and `tireAge` |
| `aggregator/engine.go` | Load stints in `setUpInitialStore`; handle `*model.Pit` in `handle`; add `getPitCounts`, `getTireCompounds`, `getTireAges`; update `GetSnapshot` |
| `datasource/historical-source.go` | Fetch stints into `raceData.Stints`; fetch pits into `eventData.EventModels` |
| `tui/model.go` | Replace hardcoded `pits`, `tires`, `tireAge` slices with snapshot data |

---

### Task 1: Add pit/tire fields to `model.Snapshot`

**Files:**
- Modify: `model/snapshot.go`

- [ ] **Step 1: Add the three new fields to `Snapshot`**

In `model/snapshot.go`, update the `Snapshot` struct:

```go
type Snapshot struct {
	SessionBar      *SessionBarSnapShot
	RaceControlMsgs []string
	Drivers         []DriverSnapshot
	LastLap         []string
	LastLapIsPitOut []bool
	Intervals       []string
	GapsToLeaders   []string
	PitCounts       []int
	TireCompounds   []string
	TireAges        []string
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./...
```

Expected: no errors. (Snapshot fields are not populated yet so no callers break.)

- [ ] **Step 3: Commit**

```bash
git add model/snapshot.go
git commit -m "Add PitCounts, TireCompounds, TireAges to Snapshot"
```

---

### Task 2: Update `Store` and `Driver` structs

**Files:**
- Modify: `aggregator/store.go`

- [ ] **Step 1: Add `Stints` to `Store` and `PitCount` to `Driver`**

In `aggregator/store.go`, update both structs:

```go
type Store struct {
	history              []model.Snapshot
	Drivers              map[uint]*Driver
	RaceControl          []model.RaceControl
	Pitstops             []model.Pit
	TotalLaps            uint
	CurrentLap           uint
	IsReplay             bool
	Session              *model.Session
	Meeting              *model.Meeting
	StartingGrid         []model.StartingGrid
	FastestLap           *FastestLap
	Stints               map[uint][]model.Stint
}

type Driver struct {
	Number           uint
	Info             *model.Driver
	Position         uint
	StartingPosition uint
	IsOut            bool
	Interval         string
	ToLeader         string
	LastLap          float64
	LastLapIsPitOut  bool
	Stint            *model.Stint
	PitCount         uint
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add aggregator/store.go
git commit -m "Add Stints map to Store and PitCount to Driver"
```

---

### Task 3: Write failing tests for tire helpers

**Files:**
- Create: `aggregator/tire_test.go`

- [ ] **Step 1: Create the test file**

```go
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
```

- [ ] **Step 2: Run tests — expect compile failure**

```bash
go test ./aggregator/...
```

Expected: compile error — `undefined: getActiveStint` and `undefined: tireAge`.

---

### Task 4: Implement tire helpers

**Files:**
- Create: `aggregator/tire.go`

- [ ] **Step 1: Create the implementation file**

```go
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
```

- [ ] **Step 2: Run tests — expect all pass**

```bash
go test ./aggregator/... -v
```

Expected: all 7 tests pass.

- [ ] **Step 3: Commit**

```bash
git add aggregator/tire.go aggregator/tire_test.go
git commit -m "Add getActiveStint and tireAge helpers with tests"
```

---

### Task 5: Fetch stints and pits in `HistoricalSource`

**Files:**
- Modify: `datasource/historical-source.go`

- [ ] **Step 1: Add stints and pits fetching to `Fetch()`**

In `HistoricalSource.Fetch()`, after the `laps` fetch block and before the assignments, add:

```go
stints, err := hs.getStints(ctx, rl, sessionKey)
if err != nil {
    return fmt.Errorf("HistoricalSource.Fetch - stints failed %w", err)
}

pits, err := hs.getPits(ctx, rl, sessionKey)
if err != nil {
    return fmt.Errorf("HistoricalSource.Fetch - pits failed %w", err)
}
```

Then in the assignment block, add `hs.raceData.Stints = stints` after the existing assignments:

```go
hs.raceData.Meeting = &meetings[0]
hs.raceData.Session = raceSession
hs.raceData.TotalLaps = getLapCount(raceControls)
hs.raceData.StartingGrid = grid
hs.raceData.Drivers = drivers
hs.raceData.SessionStart = startTime
hs.raceData.Stints = stints
```

In the event model appending block, add the pit loop after the laps loop:

```go
for _, pit := range pits {
    hs.eventData.EventModels = append(hs.eventData.EventModels, &pit)
}
```

- [ ] **Step 2: Add the two helper methods at the bottom of the file**

```go
func (hs *HistoricalSource) getStints(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.Stint, error) {
	rl.Wait()
	return hs.service.FetchStint(ctx, sessionKey)
}

func (hs *HistoricalSource) getPits(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.Pit, error) {
	rl.Wait()
	return hs.service.FetchPits(ctx, sessionKey)
}
```

- [ ] **Step 3: Verify it compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add datasource/historical-source.go
git commit -m "Fetch stints and pits from OpenF1 in HistoricalSource"
```

---

### Task 6: Wire pit/tire data into the engine

**Files:**
- Modify: `aggregator/engine.go`

- [ ] **Step 1: Load stints in `setUpInitialStore()`**

After `e.store.Drivers = convertDrivers(rd.Drivers)`, add:

```go
e.store.Stints = make(map[uint][]model.Stint)
for _, s := range rd.Stints {
    e.store.Stints[s.DriverNumber] = append(e.store.Stints[s.DriverNumber], s)
}
```

- [ ] **Step 2: Handle `*model.Pit` in `handle()`**

Add a case to the switch in `handle()`:

```go
case *model.Pit:
    e.updatePit(m)
```

- [ ] **Step 3: Add `updatePit` method**

```go
func (e *Engine) updatePit(data *model.Pit) {
	if driver, ok := e.store.Drivers[data.DriverNumber]; ok {
		driver.PitCount++
	}
}
```

- [ ] **Step 4: Add `getPitCounts`, `getTireCompounds`, `getTireAges` methods**

```go
func (e *Engine) getPitCounts() []int {
	counts := make([]int, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		counts[d.Position-1] = int(d.PitCount)
	}
	return counts
}

func (e *Engine) getTireCompounds() []string {
	compounds := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		stint := getActiveStint(e.store.Stints[d.Number], e.store.CurrentLap)
		if stint == nil {
			compounds[d.Position-1] = "---"
		} else {
			compounds[d.Position-1] = stint.Compound
		}
	}
	return compounds
}

func (e *Engine) getTireAges() []string {
	ages := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		stint := getActiveStint(e.store.Stints[d.Number], e.store.CurrentLap)
		if stint == nil {
			ages[d.Position-1] = "--"
		} else {
			ages[d.Position-1] = strconv.FormatUint(uint64(tireAge(stint, e.store.CurrentLap)), 10)
		}
	}
	return ages
}
```

- [ ] **Step 5: Update `GetSnapshot()` to populate the new fields**

In `GetSnapshot()`, update the `snapshot` literal to add the three new fields:

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
}
```

- [ ] **Step 6: Verify it compiles**

```bash
go build ./...
```

Expected: no errors. (`strconv` is already imported in `engine.go`.)

- [ ] **Step 7: Commit**

```bash
git add aggregator/engine.go
git commit -m "Wire pit count, tire compound, and tire age into engine snapshot"
```

---

### Task 7: Replace hardcoded TUI data with snapshot fields

**Files:**
- Modify: `tui/model.go`

- [ ] **Step 1: Add `strconv` to imports**

Update the import block in `tui/model.go`:

```go
import (
    "log"
    "strconv"
    "time"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "github.com/nstandage/f1-go-cli-app/aggregator"
    "github.com/nstandage/f1-go-cli-app/tui/view"
)
```

- [ ] **Step 2: Replace the three hardcoded slices in `View()`**

Remove:

```go
pits := []string{
    "1", "1", "1", "1", "0", "0", "2", "1", "0", "4",
}

tires := []string{
    "MEDIUM", "HARD", "SOFT", "MEDIUM", "MEDIUM", "SOFT", "SOFT", "INT", "WET", "SOFT",
}

tireAge := []string{
    "23", "22", "10", "17", "0", "1", "30", "29", "1", "2",
}
```

Replace with:

```go
pits := make([]string, len(snapshot.PitCounts))
for i, count := range snapshot.PitCounts {
    pits[i] = strconv.Itoa(count)
}

tires := snapshot.TireCompounds

tireAge := snapshot.TireAges
```

- [ ] **Step 3: Build and run the app to verify**

```bash
go build ./... && go run .
```

Expected: app starts, pit count shows `0` for all drivers initially then increments as pit events replay, tire compound and age update as the race progresses.

- [ ] **Step 4: Commit**

```bash
git add tui/model.go
git commit -m "Wire live pit count, tire compound, and tire age into TUI"
```
