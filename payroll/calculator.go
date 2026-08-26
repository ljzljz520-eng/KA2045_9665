package payroll

import (
	"context"
	"fmt"
	"sync"
	"time"
	"workpay/domain"
	"workpay/store"
)

type Calculator struct {
	db      *store.BoltStore
	workers int
	slots   chan struct{}
}

func New(db *store.BoltStore, workers int) *Calculator {
	if workers < 1 {
		workers = 1
	}
	return &Calculator{db: db, workers: workers, slots: make(chan struct{}, workers)}
}
func (c *Calculator) calculate(ctx context.Context, r domain.Record) (float64, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case c.slots <- struct{}{}:
	}
	// A slot has been acquired above. Release it on every return path so that
	// cancellation never strands worker capacity — the deferred drain runs
	// whether we succeed, fail, or are cancelled mid-calculation.
	defer func() { <-c.slots }()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-time.After(time.Millisecond):
	}
	return r.Amount(), nil
}
func (c *Calculator) settleOne(ctx context.Context, id string) (float64, error) {
	r, e := c.db.GetRecord(id)
	if e != nil {
		return 0, e
	}
	if ctx.Err() != nil {
		// The period is being cancelled. Do not transition the record's status
		// and let the shared ctx propagate through calculate so it can release
		// any worker resources it acquires.
		return 0, ctx.Err()
	}
	r.Status = domain.StatusProcessing
	if e = c.db.PutRecord(r); e != nil {
		return 0, e
	}
	v, e := c.calculate(ctx, r)
	if e != nil {
		// On cancellation (or any error) leave the record in its pre-settlement
		// state rather than committing a partial "processing" status.
		return 0, e
	}
	r.Status = domain.StatusSettled
	if e = c.db.PutRecord(r); e != nil {
		return 0, e
	}
	return v, nil
}
func (c *Calculator) SettlePeriod(ctx context.Context, p domain.Period) (float64, error) {
	if p.Closed {
		return 0, domain.ErrPeriodClosed
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 0.0
	errs := make(chan error, len(p.RecordIDs))
	for _, id := range p.RecordIDs {
		wg.Add(1)
		go func(recordID string) {
			defer wg.Done()
			v, e := c.settleOne(ctx, recordID)
			if e != nil {
				errs <- e
				return
			}
			mu.Lock()
			total += v
			mu.Unlock()
		}(id)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			return total, e
		}
	}
	return total, nil
}
func (c *Calculator) ResourceCount() int { return len(c.slots) }
func (c *Calculator) Describe() string {
	return fmt.Sprintf("workers=%d active=%d", c.workers, len(c.slots))
}
