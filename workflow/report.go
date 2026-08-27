package workflow

import (
	"encoding/json"
	"workpay/domain"
)

func EncodeReport(p domain.Period) []byte { b, _ := json.Marshal(p); return b }
func DescribePeriod(p domain.Period) string {
	if p.Closed {
		return "closed"
	}
	if len(p.RecordIDs) == 0 {
		return "empty"
	}
	return "open"
}
func ClonePeriod(p domain.Period) domain.Period {
	p.RecordIDs = append([]string(nil), p.RecordIDs...)
	return p
}
