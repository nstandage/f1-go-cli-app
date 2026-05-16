# Sector Colors Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Animate per-driver mini-sector color blocks in real-time during replay, progressively revealing each segment based on actual sector timing data from lap events.

**Architecture:** When a lap event arrives, a single goroutine per driver loops through all three sectors sequentially, sleeping `sectorDuration/segmentCount` between each segment reveal and appending the segment code to the driver's `Sectors` field. A `context.CancelFunc` per driver stops the previous goroutine when a new lap arrives. Sector state flows through the snapshot as `[][][]uint` (indexed by position) to the view, which renders one row of colored blocks per driver.

**Tech Stack:** Go, `context` package for goroutine cancellation, Bubbletea v2, Lipgloss v2

---

### Task 1: Extend model.Snapshot with sector fields

**Files:**
- Modify: `model/snapshot.go`

- [ ] **Step 1: Add Sectors and SectorCounts to Snapshot**

In `model/snapshot.go`, add two fields to the `Snapshot` struct:

```go
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
	Sectors         [][][]uint
	SectorCounts    [3]int
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add model/snapshot.go
git commit -m "feat: add Sectors and SectorCounts to Snapshot"
```

---

### Task 2: Extend aggregator types with sector state

**Files:**
- Modify: `aggregator/store.go`

- [ ] **Step 1: Add context import and new fields to Driver and Store**

In `aggregator/store.go`, update the import block to include `context`:

```go
import (
	"context"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)
```

Add `Sectors` and `cancelSectors` to `Driver`:

```go
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
	LapsOnTire       uint
	PitCount         uint
	CurrentLap       uint
	Sectors          [3][]uint
	cancelSectors    context.CancelFunc
}
```

Add `SectorCounts` to `Store`:

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
	SectorCounts [3]int
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add aggregator/store.go
git commit -m "feat: add sector state fields to Driver and Store"
```

---

### Task 3: Implement animateSectors and update updateLap

**Files:**
- Create: `aggregator/sector_test.go`
- Modify: `aggregator/store.go`

- [ ] **Step 1: Write failing tests**

Create `aggregator/sector_test.go`:

```go
package aggregator

import (
	"context"
	"testing"

	"github.com/nstandage/f1-go-cli-app/model"
)

func TestAnimateSectors_FillsAllSegments(t *testing.T) {
	s := &Store{}
	driver := &Driver{Number: 1}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lap := &model.Lap{
		DurationSector1: 0.003, // 3ms / 3 segments = 1ms each
		DurationSector2: 0.002, // 2ms / 2 segments = 1ms each
		DurationSector3: 0.002,
		SegmentsSector1: []uint{2049, 2049, 2049},
		SegmentsSector2: []uint{2048, 2048},
		SegmentsSector3: []uint{2051, 2051},
	}

	s.animateSectors(ctx, driver, lap)

	if len(driver.Sectors[0]) != 3 {
		t.Errorf("sector 1: expected 3 segments, got %d", len(driver.Sectors[0]))
	}
	if driver.Sectors[0][0] != 2049 {
		t.Errorf("sector 1[0]: expected 2049, got %d", driver.Sectors[0][0])
	}
	if len(driver.Sectors[1]) != 2 {
		t.Errorf("sector 2: expected 2 segments, got %d", len(driver.Sectors[1]))
	}
	if len(driver.Sectors[2]) != 2 {
		t.Errorf("sector 3: expected 2 segments, got %d", len(driver.Sectors[2]))
	}
}

func TestAnimateSectors_CancelStopsWrites(t *testing.T) {
	s := &Store{}
	driver := &Driver{Number: 1}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before calling — goroutine sleeps then sees cancelled ctx

	lap := &model.Lap{
		DurationSector1: 0.001,
		DurationSector2: 0.001,
		DurationSector3: 0.001,
		SegmentsSector1: []uint{2049},
		SegmentsSector2: []uint{2048},
		SegmentsSector3: []uint{2051},
	}

	s.animateSectors(ctx, driver, lap)

	if len(driver.Sectors[0]) != 0 {
		t.Errorf("sector 1: expected 0 segments after cancel, got %d", len(driver.Sectors[0]))
	}
	if len(driver.Sectors[1]) != 0 {
		t.Errorf("sector 2: expected 0 segments after cancel, got %d", len(driver.Sectors[1]))
	}
	if len(driver.Sectors[2]) != 0 {
		t.Errorf("sector 3: expected 0 segments after cancel, got %d", len(driver.Sectors[2]))
	}
}
```

- [ ] **Step 2: Run tests to confirm they fail**

```bash
go test ./aggregator/... -run TestAnimateSectors -v
```

Expected: compile error — `s.animateSectors undefined`

- [ ] **Step 3: Implement animateSectors**

Add to `aggregator/store.go`:

```go
func (s *Store) animateSectors(ctx context.Context, driver *Driver, data *model.Lap) {
	segments := [3][]uint{data.SegmentsSector1, data.SegmentsSector2, data.SegmentsSector3}
	durations := [3]float64{data.DurationSector1, data.DurationSector2, data.DurationSector3}

	for i, segs := range segments {
		if len(segs) == 0 {
			continue
		}
		delay := time.Duration(durations[i] / float64(len(segs)) * float64(time.Second))
		for _, seg := range segs {
			time.Sleep(delay)
			select {
			case <-ctx.Done():
				return
			default:
			}
			driver.Sectors[i] = append(driver.Sectors[i], seg)
		}
	}
}
```

- [ ] **Step 4: Update updateLap to cancel previous goroutine and launch animateSectors**

Replace `updateLap` in `aggregator/store.go` with:

```go
func (s *Store) updateLap(data *model.Lap) {
	if data.LapDuration <= 0 {
		return
	}

	if data.LapNumber > s.CurrentLap {
		s.CurrentLap = data.LapNumber
	}

	if s.SectorCounts[0] == 0 && len(data.SegmentsSector1) > 0 {
		s.SectorCounts = [3]int{
			len(data.SegmentsSector1),
			len(data.SegmentsSector2),
			len(data.SegmentsSector3),
		}
	}

	if driver, ok := s.Drivers[data.DriverNumber]; ok {
		if data.LapNumber > driver.CurrentLap {
			driver.CurrentLap = data.LapNumber
		}
		if driver.cancelSectors != nil {
			driver.cancelSectors()
		}
		driver.Sectors = [3][]uint{}
		ctx, cancel := context.WithCancel(context.Background())
		driver.cancelSectors = cancel
		go s.animateSectors(ctx, driver, data)
	}

	go s.sleepForLapDuration(data)
}
```

- [ ] **Step 5: Run tests to confirm they pass**

```bash
go test ./aggregator/... -run TestAnimateSectors -v
```

Expected:
```
--- PASS: TestAnimateSectors_FillsAllSegments
--- PASS: TestAnimateSectors_CancelStopsWrites
```

- [ ] **Step 6: Run all tests to check for regressions**

```bash
go test ./...
```

Expected: all pass.

- [ ] **Step 7: Commit**

```bash
git add aggregator/store.go aggregator/sector_test.go
git commit -m "feat: implement animateSectors goroutine with cancellation"
```

---

### Task 4: Populate sectors in GetSnapshot()

**Files:**
- Modify: `aggregator/engine.go`

- [ ] **Step 1: Add getSectors helper**

Add to `aggregator/engine.go`:

```go
func (e *Engine) getSectors() [][][]uint {
	result := make([][][]uint, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		if d.Position == 0 {
			continue
		}
		sectors := make([][]uint, 3)
		for i := range 3 {
			seg := make([]uint, len(d.Sectors[i]))
			copy(seg, d.Sectors[i])
			sectors[i] = seg
		}
		result[d.Position-1] = sectors
	}
	return result
}
```

- [ ] **Step 2: Update GetSnapshot to include sectors**

In `aggregator/engine.go`, update the `snapshot` literal in `GetSnapshot`:

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
	Sectors:         e.getSectors(),
	SectorCounts:    e.store.SectorCounts,
}
```

- [ ] **Step 3: Verify compilation and tests pass**

```bash
go build ./... && go test ./...
```

Expected: all pass.

- [ ] **Step 4: Commit**

```bash
git add aggregator/engine.go
git commit -m "feat: populate Sectors and SectorCounts in GetSnapshot"
```

---

### Task 5: Update views and wire snapshot into tui/model.go

View signature changes and the caller update must happen atomically — changing `view.Laps`'s signature breaks `tui/model.go` until the call site is updated. All changes in this task compile together before any commit.

**Files:**
- Modify: `tui/view/lap.go`
- Modify: `tui/view/top-bar.go`
- Modify: `tui/model.go`

- [ ] **Step 1: Replace Laps and miniSectorColor in tui/view/lap.go**

Replace the entire file `tui/view/lap.go` with:

```go
package view

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

func Laps(drivers [][][]uint) string {
	rows := make([]string, len(drivers))
	for i, sectors := range drivers {
		row := ""
		for _, sector := range sectors {
			for _, seg := range sector {
				row += defaultTextStyle(strings.Repeat(fullShadeBlock, 2), miniSectorColor(seg))
			}
			row += "   "
		}
		rows[i] = row
	}
	return lipgloss.NewStyle().Margin(0, 0, 0, 7).Render(
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func miniSectorColor(i uint) color.Color {
	switch i {
	case 2048:
		return slowSectorColor
	case 2049:
		return bestPersonalSectorColor
	case 2051:
		return bestOverallSectorColor
	case 2064:
		return pitLaneSectorColor
	default:
		return futureSectorColor
	}
}
```

- [ ] **Step 2: Guard sectorTitle against zero miniSectors in tui/view/top-bar.go**

In `tui/view/top-bar.go`, update `sectorTitle` to guard against a zero count (which occurs before the first lap arrives):

```go
func sectorTitle(sector string, miniSectors int) string {
	if miniSectors < 2 {
		miniSectors = 2
	}
	title := fmt.Sprintf("--%v", sector) + strings.Repeat("-", (miniSectors-2)*2)
	return lipgloss.
		NewStyle().
		Margin(0, 1).
		Render(defaultTextStyle(title, titleDarkColor))
}
```

- [ ] **Step 3: Wire snapshot into tui/model.go**

In `tui/model.go`, remove the hardcoded `lapSectors` and `lapSectorCount` variables and replace with snapshot data. Remove these lines:

```go
lapSectors := [][]int{
    {2049, 2049, 2049, 2051, 2049, 2051, 2049, 2049},
    {2049, 2049, 2049, 2049, 2049, 2049, 2049, 2049},
    {2048, 2048, 2048, 2048, 2048, 2064, 2064, 2064},
}
lapSectorCount := []int{
    len(lapSectors[0]),
    len(lapSectors[1]),
    len(lapSectors[2]),
}
topBar := view.Topbar(lapSectorCount)
```

And replace with:

```go
topBar := view.Topbar(snapshot.SectorCounts[:])
```

Then replace:

```go
laps := view.Laps(lapSectors)
```

With:

```go
laps := view.Laps(snapshot.Sectors)
```

- [ ] **Step 4: Verify compilation and all tests pass**

```bash
go build ./... && go test ./...
```

Expected: all pass.

- [ ] **Step 5: Commit**

```bash
git add tui/view/lap.go tui/view/top-bar.go tui/model.go
git commit -m "feat: render per-driver sector rows from snapshot data"
```
