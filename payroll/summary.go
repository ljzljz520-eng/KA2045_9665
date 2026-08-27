package payroll

import "workpay/domain"

func BuildSummary(records []domain.Record) map[string]float64 {
	out := map[string]float64{}
	for _, r := range records {
		out[r.WorkerID] += r.Amount()
	}
	return out
}
func Eligible(records []domain.Record) []domain.Record {
	out := make([]domain.Record, 0, len(records))
	for _, r := range records {
		if r.Status == domain.StatusSettled {
			out = append(out, r)
		}
	}
	return out
}
func Average(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	return domain.SumRecords(records) / float64(len(records))
}
