package registry

import (
	"fmt"
	"workpay/domain"
)

type BatchResult struct {
	Accepted, Rejected int
	Errors             []error
}

func (r *Registry) RegisterBatch(rs []domain.Record) BatchResult {
	out := BatchResult{}
	for _, v := range rs {
		if e := r.RegisterRecord(v); e != nil {
			out.Rejected++
			out.Errors = append(out.Errors, e)
		} else {
			out.Accepted++
		}
	}
	return out
}
func (r *Registry) ValidateBatch(rs []domain.Record, worker domain.Profile) []error {
	out := []error{}
	for _, v := range rs {
		if e := v.Validate(worker); e != nil {
			out = append(out, e)
		}
	}
	return out
}
func (r *Registry) RequireRecord(id string) (domain.Record, error) {
	v, e := r.GetRecord(id)
	if e != nil {
		return v, e
	}
	if v.ID == "" {
		return v, fmt.Errorf("record missing")
	}
	return v, nil
}
func (r *Registry) SetStatus(id, status string) error {
	v, e := r.GetRecord(id)
	if e != nil {
		return e
	}
	if !domain.Transition(v.Status, status) {
		return fmt.Errorf("invalid transition")
	}
	v.Status = status
	return r.db.PutRecord(v)
}
