# Pit Count, Tire Compound & Tire Age Feature

**Date:** 2026-05-14

## Goal

Replace the hardcoded `pits`, `tires`, and `tireAge` placeholder slices in `tui/model.go` with live data derived from the OpenF1 `pit` and `stints` API endpoints.

## Data Sources

| Column | API endpoint | Mechanism |
|--------|-------------|-----------|
| Pit count | `/pit` | Event-driven: `Pit` events replayed through event stream; `Driver.PitCount` incremented on each event |
| Tire compound | `/stints` | Reference data: fetched once, looked up by active stint in `GetSnapshot()` |
| Tire age | `/stints` | Reference data: `tyre_age_at_start + (current_lap - lap_start)` computed in `GetSnapshot()` |

## Derivations

**Active stint** — for a driver at `current_lap`: the stint with the highest `lap_start` where `lap_start <= current_lap`.

**Pit count** — incremented by 1 on each `Pit` event for the driver. Starts at 0 (no pit stops made).

**Tire compound** — `active_stint.Compound` (e.g. `"SOFT"`, `"MEDIUM"`, `"HARD"`, `"INT"`, `"WET"`).

**Tire age** — `active_stint.TyreAgeAtStart + (current_lap - active_stint.LapStart)`. Includes pre-used tyre age (e.g. qualifying tyres carried into a sprint race).

## Changes

### 1. `datasource/historical-source.go`
- Fetch stints via `service.FetchStint()`, store in `raceData.Stints`.
- Fetch pits via `service.FetchPits()`, append each `Pit` to `eventData.EventModels` (replayed in timestamp order alongside intervals/laps/positions).

### 2. `aggregator/store.go`
- Add `Stints map[uint][]model.Stint` to `Store` (keyed by `DriverNumber`).
- Add `PitCount uint` to `Driver`.

### 3. `aggregator/engine.go`
- In `setUpInitialStore()`: populate `store.Stints` from `raceData.Stints`.
- In `handle()`: add `case *model.Pit` → `store.Drivers[pit.DriverNumber].PitCount++`.
- Add `getTireSnapshot()` helper: for each driver, find active stint, return compound and tire age slices.
- In `GetSnapshot()`: add `PitCounts`, `TireCompounds`, `TireAges` fields populated from store.

### 4. `model/snapshot.go`
Add three fields to `Snapshot`:
```go
PitCounts     []int
TireCompounds []string
TireAges      []string
```

### 5. `tui/model.go`
Replace hardcoded `pits`, `tires`, `tireAge` slices with `snapshot.PitCounts` (converted to strings), `snapshot.TireCompounds`, and `snapshot.TireAges`.

## Edge Cases

- **No active stint found** (driver hasn't started yet or data gap): show `"---"` for compound, `"--"` for age, `"0"` for pit count.
- **`current_lap == 0`** (before race starts): all drivers show defaults above.
- **Driver not in stints data**: treat same as no active stint.
