package archive

import (
	"time"
	"workpay/domain"
)

func OlderThan(r domain.Record, t time.Time) bool { return r.CreatedAt.Before(t) }
func Partition(rs []domain.Record, t time.Time) ([]domain.Record, []domain.Record) {
	old, newer := []domain.Record{}, []domain.Record{}
	for _, r := range rs {
		if OlderThan(r, t) {
			old = append(old, r)
		} else {
			newer = append(newer, r)
		}
	}
	return old, newer
}
func Retain(rs []domain.Record, days int) []domain.Record {
	cut := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	_, keep := Partition(rs, cut)
	return keep
}
func ArchiveCandidates(rs []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.IsSettled() {
			out = append(out, r)
		}
	}
	return out
}
func FullyArchived(rs []domain.Record) bool {
	if len(rs) == 0 {
		return false
	}
	for _, r := range rs {
		if r.Status != domain.StatusArchived {
			return false
		}
	}
	return true
}
func ArchiveRatio(rs []domain.Record) float64 {
	if len(rs) == 0 {
		return 0
	}
	n := 0
	for _, r := range rs {
		if r.Status == domain.StatusArchived {
			n++
		}
	}
	return float64(n) / float64(len(rs))
}
