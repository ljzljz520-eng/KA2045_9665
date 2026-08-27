package payroll

import "workpay/domain"

type Analysis struct {
	Total, Minimum, Maximum, Average float64
	Count                            int
}

func Analyze(rs []domain.Record) Analysis {
	a := Analysis{Count: len(rs)}
	if len(rs) == 0 {
		return a
	}
	a.Minimum = rs[0].Amount()
	for _, r := range rs {
		v := r.Amount()
		a.Total += v
		if v < a.Minimum {
			a.Minimum = v
		}
		if v > a.Maximum {
			a.Maximum = v
		}
	}
	a.Average = a.Total / float64(a.Count)
	return a
}
func AnalyzeSettled(rs []domain.Record) Analysis { return Analyze(Eligible(rs)) }
func Taxable(a Analysis, threshold float64) float64 {
	if a.Total <= threshold {
		return 0
	}
	return a.Total - threshold
}
func Tax(a Analysis, threshold, rate float64) float64 { return Taxable(a, threshold) * rate }
func Net(a Analysis, threshold, rate float64) float64 { return a.Total - Tax(a, threshold, rate) }
func WageBands(rs []domain.Record) map[string]int {
	out := map[string]int{}
	for _, r := range rs {
		band := "low"
		if r.Amount() >= 1000 {
			band = "high"
		} else if r.Amount() >= 500 {
			band = "medium"
		}
		out[band]++
	}
	return out
}
func Bucket(rs []domain.Record, size float64) map[int]int {
	out := map[int]int{}
	if size <= 0 {
		return out
	}
	for _, r := range rs {
		out[int(r.Amount()/size)]++
	}
	return out
}
func Reconcile(expected float64, rs []domain.Record) float64 { return expected - domain.SumRecords(rs) }
func Balanced(expected float64, rs []domain.Record, tolerance float64) bool {
	d := Reconcile(expected, rs)
	if d < 0 {
		d = -d
	}
	return d <= tolerance
}
func NeedsReview(rs []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.Quantity > 1000 || r.Rate > 10000 {
			out = append(out, r)
		}
	}
	return out
}
func ValidForExport(rs []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.IsSettled() {
			out = append(out, r)
		}
	}
	return out
}
func Forecast(rs []domain.Record, days int) float64 {
	if days <= 0 {
		return 0
	}
	a := Analyze(rs)
	return a.Average * float64(days)
}
