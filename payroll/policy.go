package payroll

import "workpay/domain"

type Policy struct {
	OvertimeThreshold  float64
	OvertimeMultiplier float64
	MinimumRate        float64
}

func DefaultPolicy() Policy {
	return Policy{OvertimeThreshold: 8, OvertimeMultiplier: 1.5, MinimumRate: 1}
}
func (p Policy) Validate(r domain.Record) bool { return r.Rate >= p.MinimumRate && r.Quantity > 0 }
func (p Policy) Amount(r domain.Record) float64 {
	if r.Quantity > p.OvertimeThreshold {
		return p.OvertimeThreshold*r.Rate + (r.Quantity-p.OvertimeThreshold)*r.Rate*p.OvertimeMultiplier
	}
	return r.Amount()
}
func ApplyPolicy(rs []domain.Record, p Policy) float64 {
	total := 0.0
	for _, r := range rs {
		if p.Validate(r) {
			total += p.Amount(r)
		}
	}
	return total
}
