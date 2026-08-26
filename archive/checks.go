package archive

import "workpay/domain"

func CheckIDs(ids []string) bool {
	return domain.EnsureUnique(ids)
}
func CheckAmounts(rs []domain.Record) bool {
	for _, r := range rs {
		if r.Amount() < 0 {
			return false
		}
	}
	return true
}
func CheckStatuses(rs []domain.Record, status string) bool {
	for _, r := range rs {
		if r.Status != status {
			return false
		}
	}
	return true
}
func CountArchived(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		if r.Status == domain.StatusArchived {
			n++
		}
	}
	return n
}
func CountSettled(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		if r.Status == domain.StatusSettled {
			n++
		}
	}
	return n
}
func CountPending(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		if r.Status == domain.StatusPending {
			n++
		}
	}
	return n
}
func CountProcessing(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		if r.Status == domain.StatusProcessing {
			n++
		}
	}
	return n
}
func TotalArchived(rs []domain.Record) float64 {
	total := 0.0
	for _, r := range rs {
		if r.Status == domain.StatusArchived {
			total += r.Amount()
		}
	}
	return total
}
func TotalSettled(rs []domain.Record) float64 {
	total := 0.0
	for _, r := range rs {
		if r.Status == domain.StatusSettled {
			total += r.Amount()
		}
	}
	return total
}
func SameWorker(rs []domain.Record, id string) bool {
	for _, r := range rs {
		if r.WorkerID != id {
			return false
		}
	}
	return true
}
func SameSite(rs []domain.Record, id string) bool {
	for _, r := range rs {
		if r.SiteID != id {
			return false
		}
	}
	return true
}
