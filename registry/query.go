package registry

import (
	"encoding/json"
	"workpay/domain"
)

func (r *Registry) ListRecords() ([]domain.Record, error) {
	raw, e := r.db.List("records")
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, b := range raw {
		var v domain.Record
		if json.Unmarshal(b, &v) == nil {
			out = append(out, v)
		}
	}
	return out, nil
}
func (r *Registry) ListProfiles() ([]domain.Profile, error) {
	raw, e := r.db.List("profiles")
	if e != nil {
		return nil, e
	}
	out := []domain.Profile{}
	for _, b := range raw {
		var v domain.Profile
		if json.Unmarshal(b, &v) == nil {
			out = append(out, v)
		}
	}
	return out, nil
}
