package archive

import (
	"fmt"
	"workpay/domain"
)

func FormatPeriod(p domain.Period) string {
	state := "open"
	if p.Closed {
		state = "closed"
	}
	return fmt.Sprintf("period=%s state=%s total=%.2f records=%d", p.ID, state, p.Total, len(p.RecordIDs))
}
func CanArchive(p domain.Period, records []domain.Record) bool {
	if p.Closed || len(records) != len(p.RecordIDs) {
		return false
	}
	for _, r := range records {
		if !r.IsSettled() {
			return false
		}
	}
	return true
}
