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
