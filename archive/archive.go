package archive

import (
	"fmt"
	"workpay/domain"
	"workpay/store"
)

type Service struct{ db *store.BoltStore }

func New(db *store.BoltStore) *Service { return &Service{db: db} }
func (s *Service) ArchivePeriod(p domain.Period, actor string) (domain.Period, error) {
	if p.Closed {
		return p, domain.ErrPeriodClosed
	}
	var total float64
	for _, id := range p.RecordIDs {
		r, e := s.db.GetRecord(id)
		if e != nil {
			return p, e
		}
		if !r.IsSettled() {
			return p, fmt.Errorf("record %s not settled", id)
		}
		r = s.archiveRecord(r)
		total += r.Amount()
		if e = s.db.PutRecord(r); e != nil {
			return p, e
		}
	}
	p = p.Close(total)
	if e := s.db.PutPeriod(p); e != nil {
		return p, e
	}
	if e := s.db.PutAudit(domain.NewAudit("audit-"+p.ID, actor, "archive", p.ID)); e != nil {
		return p, e
	}
	return p, nil
}
func (s *Service) archiveRecord(r domain.Record) domain.Record { return r.MarkArchived() }
func (s *Service) Audit(actor, action, target string) error {
	if actor == "" || action == "" || target == "" {
		return fmt.Errorf("audit fields required")
	}
	return s.db.PutAudit(domain.NewAudit(fmt.Sprintf("audit-%s-%s", actor, target), actor, action, target))
}
