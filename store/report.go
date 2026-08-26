package store

import (
	"encoding/json"
	"sort"
	"workpay/domain"
)

func DecodeRecords(raw [][]byte) []domain.Record {
	out := []domain.Record{}
	for _, b := range raw {
		var r domain.Record
		if json.Unmarshal(b, &r) == nil {
			out = append(out, r)
		}
	}
	return out
}
func DecodeProfiles(raw [][]byte) []domain.Profile {
	out := []domain.Profile{}
	for _, b := range raw {
		var p domain.Profile
		if json.Unmarshal(b, &p) == nil {
			out = append(out, p)
		}
	}
	return out
}
func DecodeEvents(raw [][]byte) []domain.Event {
	out := []domain.Event{}
	for _, b := range raw {
		var e domain.Event
		if json.Unmarshal(b, &e) == nil {
			out = append(out, e)
		}
	}
	return out
}
func DecodeAudits(raw [][]byte) []domain.Audit {
	out := []domain.Audit{}
	for _, b := range raw {
		var a domain.Audit
		if json.Unmarshal(b, &a) == nil {
			out = append(out, a)
		}
	}
	return out
}
func SortRecordsByID(rs []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rs...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func SortProfilesByName(ps []domain.Profile) []domain.Profile {
	out := append([]domain.Profile(nil), ps...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func SortEventsByType(es []domain.Event) []domain.Event {
	out := append([]domain.Event(nil), es...)
	sort.Slice(out, func(i, j int) bool { return out[i].Type < out[j].Type })
	return out
}
func SortAuditsByActor(as []domain.Audit) []domain.Audit {
	out := append([]domain.Audit(nil), as...)
	sort.Slice(out, func(i, j int) bool { return out[i].Actor < out[j].Actor })
	return out
}
func RecordsForWorker(rs []domain.Record, id string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.WorkerID == id {
			out = append(out, r)
		}
	}
	return out
}
func RecordsForSite(rs []domain.Record, id string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.SiteID == id {
			out = append(out, r)
		}
	}
	return out
}
func EventsForRecord(es []domain.Event, id string) []domain.Event {
	out := []domain.Event{}
	for _, e := range es {
		if e.RecordID == id {
			out = append(out, e)
		}
	}
	return out
}
func AuditsForTarget(as []domain.Audit, id string) []domain.Audit {
	out := []domain.Audit{}
	for _, a := range as {
		if a.Target == id {
			out = append(out, a)
		}
	}
	return out
}
func CountRecords(rs []domain.Record) int                 { return len(rs) }
func CountProfiles(ps []domain.Profile) int               { return len(ps) }
func CountEvents(es []domain.Event) int                   { return len(es) }
func CountAudits(as []domain.Audit) int                   { return len(as) }
func MarshalRecords(rs []domain.Record) ([]byte, error)   { return json.Marshal(rs) }
func MarshalProfiles(ps []domain.Profile) ([]byte, error) { return json.Marshal(ps) }
func MarshalEvents(es []domain.Event) ([]byte, error)     { return json.Marshal(es) }
func MarshalAudits(as []domain.Audit) ([]byte, error)     { return json.Marshal(as) }
func UnmarshalRecords(b []byte) ([]domain.Record, error) {
	var v []domain.Record
	e := json.Unmarshal(b, &v)
	return v, e
}
func UnmarshalProfiles(b []byte) ([]domain.Profile, error) {
	var v []domain.Profile
	e := json.Unmarshal(b, &v)
	return v, e
}
func UnmarshalEvents(b []byte) ([]domain.Event, error) {
	var v []domain.Event
	e := json.Unmarshal(b, &v)
	return v, e
}
func UnmarshalAudits(b []byte) ([]domain.Audit, error) {
	var v []domain.Audit
	e := json.Unmarshal(b, &v)
	return v, e
}
