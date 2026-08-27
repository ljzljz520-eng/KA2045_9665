package domain

import "fmt"

func ValidateRecordIDs(p Period, rs []Record) error {
	if len(p.RecordIDs) != len(rs) {
		return fmt.Errorf("record count mismatch")
	}
	seen := map[string]bool{}
	for _, r := range rs {
		if seen[r.ID] {
			return fmt.Errorf("duplicate record %s", r.ID)
		}
		seen[r.ID] = true
		if !ContainsRecord(p, r.ID) {
			return fmt.Errorf("record %s outside period", r.ID)
		}
	}
	return nil
}
func CanProcess(p Period) bool              { return p.ID != "" && !p.Closed && len(p.RecordIDs) > 0 }
func CanClose(p Period, total float64) bool { return CanProcess(p) && total >= 0 }
func NextStatus(s string) string {
	switch s {
	case StatusPending:
		return StatusProcessing
	case StatusProcessing:
		return StatusSettled
	case StatusSettled:
		return StatusArchived
	default:
		return s
	}
}
func IsTerminal(s string) bool { return s == StatusArchived }
func IsMutable(s string) bool  { return s == StatusPending || s == StatusProcessing }
func TransitionPath(from, to string) []string {
	out := []string{from}
	for from != to {
		next := NextStatus(from)
		if next == from {
			break
		}
		out = append(out, next)
		from = next
	}
	return out
}
func PeriodStatus(p Period) string {
	if p.Closed {
		return "closed"
	}
	if len(p.RecordIDs) == 0 {
		return "empty"
	}
	return "open"
}
func PeriodTotal(rs []Record) float64 {
	total := 0.0
	for _, r := range rs {
		if r.Status != StatusArchived {
			total += r.Amount()
		}
	}
	return total
}
func Archivable(rs []Record) bool {
	if len(rs) == 0 {
		return false
	}
	for _, r := range rs {
		if !r.IsSettled() {
			return false
		}
	}
	return true
}
func ValidateArchive(p Period, rs []Record) error {
	if p.Closed {
		return ErrPeriodClosed
	}
	if e := ValidateRecordIDs(p, rs); e != nil {
		return e
	}
	if !Archivable(rs) {
		return fmt.Errorf("records not settled")
	}
	return nil
}
