package payroll

import "workpay/domain"

type CostBreakdown struct{ Labor, Overtime, Allowance, Total float64 }

func Breakdown(rs []domain.Record, p Policy, allowance float64) CostBreakdown {
	b := CostBreakdown{Allowance: allowance}
	for _, r := range rs {
		base := r.Amount()
		b.Labor += base
		premium := p.Amount(r) - base
		if premium > 0 {
			b.Overtime += premium
		}
	}
	b.Total = b.Labor + b.Overtime + b.Allowance
	return b
}
func AddAllowance(b CostBreakdown, v float64) CostBreakdown {
	b.Allowance += v
	b.Total = b.Labor + b.Overtime + b.Allowance
	return b
}
func WithoutAllowance(b CostBreakdown) CostBreakdown {
	b.Total = b.Labor + b.Overtime
	b.Allowance = 0
	return b
}
func CostPerRecord(b CostBreakdown, n int) float64 {
	if n <= 0 {
		return 0
	}
	return b.Total / float64(n)
}
func CostPerWorker(b CostBreakdown, workers int) float64 {
	if workers <= 0 {
		return 0
	}
	return b.Total / float64(workers)
}
func CompareCosts(a, b CostBreakdown) float64 { return a.Total - b.Total }
func CostWithin(a, b CostBreakdown, tolerance float64) bool {
	d := CompareCosts(a, b)
	if d < 0 {
		d = -d
	}
	return d <= tolerance
}
func TaxBreakdown(b CostBreakdown, rate float64) float64 {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return b.Total * rate
}
func NetBreakdown(b CostBreakdown, rate float64) CostBreakdown {
	b.Total -= TaxBreakdown(b, rate)
	return b
}
func RoundCents(v float64) float64 { return float64(int(v*100+0.5)) / 100 }
