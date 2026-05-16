# Sector Colors Design

**Date:** 2026-05-15

## Overview

Animate per-driver sector mini-segment colors in real-time during replay, using lap sector duration data to progressively reveal each mini-segment block at the correct moment.

## Data Model

### `aggregator.Driver` (new fields)

- `Sectors [3][]uint` — currently visible mini-segments per sector, filled progressively by the sector goroutine. Empty until first lap arrives.
- `cancelSectors context.CancelFunc` — cancels the running sector goroutine for this driver. Nil until first lap.

### `aggregator.Store` (new field)

- `SectorCounts [3]int` — mini-segment counts per sector (e.g. `[8, 8, 8]`). Set once, on the first lap event that has non-empty segment data (`SectorCounts[0] == 0 && len(data.SegmentsSector1) > 0`). Never updated again — track layout is constant within a session.

### `model.Snapshot` (new fields)

- `Sectors [][][]uint` — indexed by driver position (same ordering as `Drivers`, `LastLap`, etc.). Empty outer slice until first lap arrives.
- `SectorCounts [3]int` — zero value before first lap; view renders nothing in that case.

## Goroutine Logic

In `store.updateLap`, after the existing `go s.sleepForLapDuration(data)` call:

1. If `driver.cancelSectors != nil`, call it to cancel the previous sector goroutine.
2. Reset `driver.Sectors` to `[3][]uint{}`.
3. Create `ctx, cancel := context.WithCancel(context.Background())`, store `cancel` on `driver.cancelSectors`.
4. Launch `go s.animateSectors(ctx, driver, data)`.

`animateSectors` loops through sectors sequentially:

```
for each sector N in [0, 1, 2]:
    segments := SegmentsSectorN from data
    if len(segments) == 0: continue
    delay := DurationSectorN (seconds) / len(segments), converted to time.Duration
    for each segment in segments:
        sleep(delay)
        if ctx.Done(): return   // check before writing
        append segment to driver.Sectors[N]
```

Checking `ctx.Done()` before appending ensures a cancelled goroutine never writes a stale segment.

## Snapshot Population

In `Engine.GetSnapshot()`:

- **`Sectors`**: loop over drivers by position (same pattern as `getLastLap`), copy `driver.Sectors`. Drivers with no data yet contribute `[][]uint{}`.
- **`SectorCounts`**: copy `store.SectorCounts` directly.

## View Changes

### `tui/view/lap.go`

- `Laps` signature changes from `Laps(sectors [][]int)` to `Laps(drivers [][][]uint)`.
- Outer loop produces one string per driver row using existing sector/mini-segment block rendering.
- Rows are joined with `lipgloss.JoinVertical` and wrapped in the existing margin style.
- `miniSectorColor` parameter changes from `int` to `uint` to match model type.

### `tui/model.go`

- Remove hardcoded `lapSectors` and `lapSectorCount` variables.
- Pass `snapshot.Sectors` to `view.Laps`.
- Pass `snapshot.SectorCounts[:]` to `view.Topbar`.

## Initial State

Sectors render blank (nothing) until the first lap with non-empty segment data arrives. No placeholder blocks are shown.
