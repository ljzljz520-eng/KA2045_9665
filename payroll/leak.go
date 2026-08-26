package payroll

import (
	"context"
	"workpay/domain"
)

// RunWithReservation settles a single record while respecting cancellation.
//
// Earlier this spawned a goroutine that acquired a worker slot and then called
// calculate with a background (uncancellable) context, while the caller
// returned immediately on ctx.Done(). That orphaned goroutine held the slot
// forever — the symptom was worker capacity leaking away after a period was
// cancelled/closed. Now the work runs in-line under the supplied ctx so a
// cancellation unwinds it promptly and releases every resource it took.
func (c *Calculator) RunWithReservation(ctx context.Context, r domain.Record) error {
	_, err := c.calculate(ctx, r)
	return err
}
