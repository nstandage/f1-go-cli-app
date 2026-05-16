package aggregator

import (
	"context"
	"testing"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

func TestAnimateSectors_FillsAllSegments(t *testing.T) {
	s := &Store{}
	driver := &Driver{
		Number: 1,
		Sectors: [3][]uint{
			make([]uint, 3),
			make([]uint, 2),
			make([]uint, 2),
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lap := &model.Lap{
		DurationSector1: 0.003,
		DurationSector2: 0.002,
		DurationSector3: 0.002,
		SegmentsSector1: []uint{2049, 2049, 2049},
		SegmentsSector2: []uint{2048, 2048},
		SegmentsSector3: []uint{2051, 2051},
	}

	s.animateSectors(ctx, 0, driver, lap)

	for idx, want := range []uint{2049, 2049, 2049} {
		if driver.Sectors[0][idx] != want {
			t.Errorf("sector 1[%d]: expected %d, got %d", idx, want, driver.Sectors[0][idx])
		}
	}
	for idx, want := range []uint{2048, 2048} {
		if driver.Sectors[1][idx] != want {
			t.Errorf("sector 2[%d]: expected %d, got %d", idx, want, driver.Sectors[1][idx])
		}
	}
	for idx, want := range []uint{2051, 2051} {
		if driver.Sectors[2][idx] != want {
			t.Errorf("sector 3[%d]: expected %d, got %d", idx, want, driver.Sectors[2][idx])
		}
	}
}

func TestAnimateSectors_CancelStopsWrites(t *testing.T) {
	s := &Store{}
	driver := &Driver{
		Number: 1,
		Sectors: [3][]uint{
			make([]uint, 1),
			make([]uint, 1),
			make([]uint, 1),
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before calling — function sleeps then sees cancelled ctx on first check

	lap := &model.Lap{
		DurationSector1: 0.001,
		DurationSector2: 0.001,
		DurationSector3: 0.001,
		SegmentsSector1: []uint{2049},
		SegmentsSector2: []uint{2048},
		SegmentsSector3: []uint{2051},
	}

	s.animateSectors(ctx, 0, driver, lap)

	if driver.Sectors[0][0] != 0 {
		t.Errorf("sector 1[0]: expected 0 (future) after cancel, got %d", driver.Sectors[0][0])
	}
	if driver.Sectors[1][0] != 0 {
		t.Errorf("sector 2[0]: expected 0 (future) after cancel, got %d", driver.Sectors[1][0])
	}
	if driver.Sectors[2][0] != 0 {
		t.Errorf("sector 3[0]: expected 0 (future) after cancel, got %d", driver.Sectors[2][0])
	}
}

func TestAnimateSectors_CancelMidWay(t *testing.T) {
	s := &Store{}
	driver := &Driver{
		Number: 1,
		Sectors: [3][]uint{
			make([]uint, 3),
			make([]uint, 3),
			make([]uint, 3),
		},
	}
	ctx, cancel := context.WithCancel(context.Background())

	lap := &model.Lap{
		DurationSector1: 0.003,
		DurationSector2: 0.003,
		DurationSector3: 0.003,
		SegmentsSector1: []uint{2049, 2049, 2049},
		SegmentsSector2: []uint{2048, 2048, 2048},
		SegmentsSector3: []uint{2051, 2051, 2051},
	}

	done := make(chan struct{})
	go func() {
		s.animateSectors(ctx, 0, driver, lap)
		close(done)
	}()

	time.Sleep(5 * time.Millisecond)
	cancel()
	<-done

	// Sector 1 should be fully written
	for idx := range 3 {
		if driver.Sectors[0][idx] != 2049 {
			t.Errorf("sector 1[%d]: expected 2049, got %d", idx, driver.Sectors[0][idx])
		}
	}
	// Sector 3 should be untouched (all future/zero)
	for idx := range 3 {
		if driver.Sectors[2][idx] != 0 {
			t.Errorf("sector 3[%d]: expected 0 (future), got %d", idx, driver.Sectors[2][idx])
		}
	}
}
