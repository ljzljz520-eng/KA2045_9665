package payroll

import (
	"context"
	"sync"
	"workpay/domain"
)

type WorkerPool struct {
	size    int
	jobs    chan domain.Record
	results chan float64
	wg      sync.WaitGroup
}

func NewWorkerPool(size int) *WorkerPool {
	if size < 1 {
		size = 1
	}
	return &WorkerPool{size: size, jobs: make(chan domain.Record), results: make(chan float64)}
}
func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.size; i++ {
		p.wg.Add(1)
		go p.run(ctx)
	}
}
func (p *WorkerPool) run(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case r, ok := <-p.jobs:
			if !ok {
				return
			}
			select {
			case p.results <- r.Amount():
			case <-ctx.Done():
				return
			}
		}
	}
}
func (p *WorkerPool) Submit(ctx context.Context, r domain.Record) error {
	select {
	case p.jobs <- r:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (p *WorkerPool) Results() <-chan float64 { return p.results }
func (p *WorkerPool) Stop()                   { close(p.jobs); p.wg.Wait(); close(p.results) }
func (p *WorkerPool) Size() int               { return p.size }
func Collect(ctx context.Context, results <-chan float64, n int) float64 {
	total := 0.0
	for i := 0; i < n; i++ {
		select {
		case v, ok := <-results:
			if !ok {
				return total
			}
			total += v
		case <-ctx.Done():
			return total
		}
	}
	return total
}
