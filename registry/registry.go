package registry

import (
	"fmt"
	"sync"
	"workpay/domain"
	"workpay/store"
)

type Registry struct {
	db *store.BoltStore
	mu sync.Mutex
}

func New(db *store.BoltStore) *Registry { return &Registry{db: db} }
func (r *Registry) RegisterProfile(p domain.Profile) error {
	if e := domain.ValidateProfile(p); e != nil {
		return e
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.db.PutProfile(p)
}
func (r *Registry) RegisterRecord(rec domain.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, e := r.db.GetProfile(rec.WorkerID)
	if e != nil {
		return e
	}
	if e = rec.Validate(p); e != nil {
		return e
	}
	if _, e = r.db.GetRecord(rec.ID); e == nil {
		return fmt.Errorf("record exists")
	}
	if e = r.db.PutRecord(rec); e != nil {
		return e
	}
	return r.db.PutEvent(domain.NewEvent("event-"+rec.ID, rec.ID, "registered", "record accepted"))
}
func (r *Registry) CreatePeriod(id string, ids []string) (domain.Period, error) {
	p := domain.Period{ID: id}
	for _, id := range ids {
		p = p.Add(id)
	}
	if e := domain.ValidatePeriod(p); e != nil {
		return p, e
	}
	if !domain.EnsureUnique(ids) {
		return p, fmt.Errorf("duplicate records")
	}
	e := r.db.PutPeriod(p)
	return p, e
}
func (r *Registry) GetRecord(id string) (domain.Record, error) { return r.db.GetRecord(id) }
