package store

import (
	"fmt"
	"time"
	"workpay/domain"
)

func (s *BoltStore) SaveAll(rs []domain.Record) error {
	for _, r := range rs {
		if e := s.PutRecord(r); e != nil {
			return e
		}
	}
	return nil
}
func (s *BoltStore) SaveProfiles(ps []domain.Profile) error {
	for _, p := range ps {
		if e := s.PutProfile(p); e != nil {
			return e
		}
	}
	return nil
}
func (s *BoltStore) SaveEvents(es []domain.Event) error {
	for _, e := range es {
		if e.CreatedAt.IsZero() {
			e.CreatedAt = time.Now().UTC()
		}
		if x := s.PutEvent(e); x != nil {
			return x
		}
	}
	return nil
}
func (s *BoltStore) SaveAudits(as []domain.Audit) error {
	for _, a := range as {
		if x := s.PutAudit(a); x != nil {
			return x
		}
	}
	return nil
}
func (s *BoltStore) MustRecord(id string) (domain.Record, error) {
	r, e := s.GetRecord(id)
	if e != nil {
		return r, e
	}
	if r.ID == "" {
		return r, fmt.Errorf("record %s missing", id)
	}
	return r, nil
}
func (s *BoltStore) HasRecord(id string) bool  { _, e := s.GetRecord(id); return e == nil }
func (s *BoltStore) HasProfile(id string) bool { _, e := s.GetProfile(id); return e == nil }
func (s *BoltStore) HasPeriod(id string) bool  { _, e := s.GetPeriod(id); return e == nil }
