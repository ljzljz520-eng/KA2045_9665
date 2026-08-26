package store

import (
	"encoding/json"
	"workpay/domain"
)

func EncodeRecord(r domain.Record) ([]byte, error) { return json.Marshal(r) }
func DecodeRecord(b []byte) (domain.Record, error) {
	var r domain.Record
	e := json.Unmarshal(b, &r)
	return r, e
}
func EncodeProfile(p domain.Profile) ([]byte, error) { return json.Marshal(p) }
func DecodeProfile(b []byte) (domain.Profile, error) {
	var p domain.Profile
	e := json.Unmarshal(b, &p)
	return p, e
}
func EncodeEvent(v domain.Event) ([]byte, error) { return json.Marshal(v) }
func EncodeAudit(v domain.Audit) ([]byte, error) { return json.Marshal(v) }
