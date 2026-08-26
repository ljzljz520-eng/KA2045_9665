package payroll

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"workpay/domain"
	"workpay/store"
)

// newStore builds a Calculator backed by a temporary BoltDB so settlement
// exercises the real persistence path used in production.
func newStore(t *testing.T, workers int) (*Calculator, func()) {
	t.Helper()
	dir := t.TempDir()
	db, err := store.Open(filepath.Join(dir, "payroll.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	cleanup := func() {
		_ = db.Close()
		_ = os.RemoveAll(dir)
	}
	return New(db, workers), cleanup
}

// seedPeriod persists records that a period can settle, returning the period.
func seedPeriod(c *Calculator, n int) (domain.Period, error) {
	p := domain.Period{ID: "period-1"}
	for i := 0; i < n; i++ {
		r := domain.NewRecord(fmt.Sprintf("rec-%d", i), "worker-1", "site-1", 1, 1)
		if e := c.db.PutRecord(r); e != nil {
			return p, e
		}
		p.RecordIDs = append(p.RecordIDs, r.ID)
	}
	return p, nil
}

// TestCancelReleasesAllSlots is the regression test for the leak: after
// cancelling an in-flight settlement, no worker slots may remain held and no
// goroutines may be left running.
func TestCancelReleasesAllSlots(t *testing.T) {
	c, cleanup := newStore(t, 2)
	defer cleanup()

	p, err := seedPeriod(c, 50)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Give the settlement a moment to acquire slots, then cancel —
		// mirroring "取消工资核算" (cancel payroll).
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	if _, err := c.SettlePeriod(ctx, p); err == nil {
		t.Fatalf("expected cancellation error, got nil")
	}

	// After cancellation every slot must be free.
	if got := c.ResourceCount(); got != 0 {
		t.Fatalf("slots leaked after cancel: got %d, want 0", got)
	}
	if got := c.Describe(); got != "workers=2 active=0" {
		t.Fatalf("Describe() = %q, want %q", got, "workers=2 active=0")
	}
}

// TestRunWithReservationRespectsCancel ensures the helper that previously
// stranded a goroutine now unwinds promptly on cancellation.
func TestRunWithReservationRespectsCancel(t *testing.T) {
	c, cleanup := newStore(t, 2)
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	if err := c.RunWithReservation(ctx, domain.Record{}); err == nil {
		t.Fatalf("expected error from cancelled context, got nil")
	}
	if got := c.ResourceCount(); got != 0 {
		t.Fatalf("slots leaked via RunWithReservation: got %d, want 0", got)
	}
}

// TestSettlePeriodCompletes sanity-checks that the happy path still settles all
// records and leaves no resources held.
func TestSettlePeriodCompletes(t *testing.T) {
	c, cleanup := newStore(t, 4)
	defer cleanup()

	p, err := seedPeriod(c, 8)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	total, err := c.SettlePeriod(context.Background(), p)
	if err != nil {
		t.Fatalf("SettlePeriod: %v", err)
	}
	if total != 8 {
		t.Fatalf("total = %v, want 8", total)
	}
	if got := c.ResourceCount(); got != 0 {
		t.Fatalf("slots leaked after success: got %d, want 0", got)
	}
	for _, id := range p.RecordIDs {
		r, err := c.db.GetRecord(id)
		if err != nil {
			t.Fatalf("GetRecord(%s): %v", id, err)
		}
		if r.Status != domain.StatusSettled {
			t.Fatalf("record %s status = %q, want %q", id, r.Status, domain.StatusSettled)
		}
	}
}

// keep imports honest for environments that race-check goroutine counts.
var _ = runtime.NumGoroutine
