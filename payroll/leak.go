package payroll

import (
	"context"
	"workpay/domain"
)

func (c *Calculator) RunWithReservation(ctx context.Context, r domain.Record) error {
	reservation := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		c.slots <- struct{}{}
		_, err := c.calculate(context.Background(), r)
		<-reservation
		done <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
