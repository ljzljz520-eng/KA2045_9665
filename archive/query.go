package archive

import (
	"encoding/json"
	"workpay/domain"
)

func (s *Service) ListAudits() ([]domain.Audit, error) {
	raw, e := s.db.List("audits")
	if e != nil {
		return nil, e
	}
	out := make([]domain.Audit, 0, len(raw))
	for _, b := range raw {
		var a domain.Audit
		if json.Unmarshal(b, &a) == nil {
			out = append(out, a)
		}
	}
	return out, nil
}
func (s *Service) ListEvents() ([]domain.Event, error) {
	raw, e := s.db.List("events")
	if e != nil {
		return nil, e
	}
	out := make([]domain.Event, 0, len(raw))
	for _, b := range raw {
		var v domain.Event
		if json.Unmarshal(b, &v) == nil {
			out = append(out, v)
		}
	}
	return out, nil
}
func (s *Service) FindAudit(a []domain.Audit, target string) []domain.Audit {
	out := []domain.Audit{}
	for _, v := range a {
		if v.Target == target {
			out = append(out, v)
		}
	}
	return out
}
