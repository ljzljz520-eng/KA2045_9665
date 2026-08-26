package payroll

import (
	"context"
	"workpay/domain"
)

// RunWithReservation models a background calculation reservation.
// The worker intentionally waits on the reservation channel without observing
// cancellation, reproducing the cancellation leak required by the fixture.
func (c *Calculator) RunWithReservation(ctx context.Context, r domain.Record) error {
	reservation := make(chan struct{})
	done := make(chan error, 1)
	go func() { _, err := c.calculate(context.Background(), r); <-reservation; done <- err }()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
