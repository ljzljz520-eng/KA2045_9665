package payroll

import (
	"sort"
	"workpay/domain"
)

type LedgerEntry struct {
	RecordID, WorkerID string
	Amount             float64
}

func Entries(rs []domain.Record) []LedgerEntry {
	out := make([]LedgerEntry, 0, len(rs))
	for _, r := range rs {
		out = append(out, LedgerEntry{r.ID, r.WorkerID, r.Amount()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Amount > out[j].Amount })
	return out
}
func TopWorkers(rs []domain.Record, n int) []LedgerEntry {
	e := Entries(rs)
	if n < 0 {
		n = 0
	}
	if n > len(e) {
		n = len(e)
	}
	return e[:n]
}
func WorkerTotals(rs []domain.Record) map[string]float64 {
	out := map[string]float64{}
	for _, r := range rs {
		out[r.WorkerID] += r.Amount()
	}
	return out
}
