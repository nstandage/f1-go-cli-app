# Live Pit Stop Feed — Design Spec

**Date:** 2026-05-15
**Branch:** pitstop-time

## Overview

Replace the static hardcoded pit stop list in the TUI with a live feed of the last 8 pit stops, showing driver acronym and stop duration for each entry. Newest stops appear at the top.

## Data Model

Add to `model/snapshot.go`:

```go
type PitStopEntry struct {
    DriverAcronym string
    StopDuration  float64
}
```

Add `RecentPitStops []PitStopEntry` to `model.Snapshot`.

## Aggregator — Store

Add `RecentPits []PitStopEntry` to `aggregator.Store` (max 8 entries, newest first).

Update `store.updatePit(p *model.Pit)` to:
1. Look up `NameAcronym` via `s.Drivers[p.DriverNumber].Info.NameAcronym`.
2. Prepend a new `PitStopEntry{DriverAcronym, StopDuration}` to `RecentPits`.
3. Trim the slice to 8 entries if it exceeds that length.

The existing `driver.PitCount++` logic is unchanged.

## Aggregator — Snapshot

In `Engine.GetSnapshot()`, copy `e.store.RecentPits` into `snapshot.RecentPitStops`.

## View

Update `tui/view/pitstop.go`:
- Change signature from `PitStops(stops []float64)` to `PitStops(stops []model.PitStopEntry)`.
- Render each row as `"VER  2.95"` — acronym left-padded to a fixed width, stop time formatted `%.2f`.
- Retain existing color logic: green (`< 3.1`), yellow (`< 3.5`), red (`>= 3.5`).
- Border width and height unchanged.

## TUI Wiring

In `tui/model.go`:
- Remove hardcoded `pitStops := []float64{...}`.
- Pass `snapshot.RecentPitStops` directly to `view.PitStops(...)`.

## What is not changing

- The `Pit` model struct (`model/pit.go`) — already has all needed fields.
- Pit event handling in `engine.handle()` — already routes `*model.Pit` to `updatePit`.
- Data fetching in `datasource/historical-source.go` — pits are already fetched and replayed.
- The border dimensions of the pit stop panel.
