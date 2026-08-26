package archive

import (
	"encoding/json"
	"workpay/domain"
)

type Export struct {
	Period  domain.Period
	Records []domain.Record
	Audits  []domain.Audit
}

func BuildExport(p domain.Period, rs []domain.Record, as []domain.Audit) Export {
	return Export{p, append([]domain.Record(nil), rs...), append([]domain.Audit(nil), as...)}
}
func MarshalExport(e Export) ([]byte, error) { return json.MarshalIndent(e, "", "  ") }
func RecordCount(e Export) int               { return len(e.Records) }
