package domain

import (
	"sort"
	"strings"
	"time"
)

type RuleResult struct {
	Allowed bool
	Code    string
	Message string
}

func CheckIdentity(id string) RuleResult {
	id = strings.TrimSpace(id)
	if id == "" {
		return RuleResult{false, "missing_id", "identity is required"}
	}
	if len(id) > 64 {
		return RuleResult{false, "long_id", "identity is too long"}
	}
	return RuleResult{true, "ok", ""}
}
func CheckQuantity(q float64) RuleResult {
	if q <= 0 {
		return RuleResult{false, "quantity", "quantity must be positive"}
	}
	if q > 100000 {
		return RuleResult{false, "quantity_limit", "quantity exceeds limit"}
	}
	return RuleResult{true, "ok", ""}
}
func CheckRate(rate float64) RuleResult {
	if rate <= 0 {
		return RuleResult{false, "rate", "rate must be positive"}
	}
	if rate > 1000000 {
		return RuleResult{false, "rate_limit", "rate exceeds limit"}
	}
	return RuleResult{true, "ok", ""}
}
func CheckTimestamp(t time.Time) RuleResult {
	if t.IsZero() {
		return RuleResult{false, "timestamp", "timestamp is required"}
	}
	if t.After(time.Now().Add(time.Hour)) {
		return RuleResult{false, "future", "timestamp is too far ahead"}
	}
	return RuleResult{true, "ok", ""}
}
func ValidateRecordDeep(r Record, p Profile) []RuleResult {
	return []RuleResult{CheckIdentity(r.ID), CheckIdentity(r.WorkerID), CheckIdentity(r.SiteID), CheckQuantity(r.Quantity), CheckRate(r.Rate), CheckTimestamp(r.CreatedAt), profileRule(p)}
}
func profileRule(p Profile) RuleResult {
	if !p.Active {
		return RuleResult{false, "inactive", "profile inactive"}
	}
	return RuleResult{true, "ok", ""}
}
func AllAllowed(rs []RuleResult) bool {
	for _, r := range rs {
		if !r.Allowed {
			return false
		}
	}
	return true
}
func FirstFailure(rs []RuleResult) (RuleResult, bool) {
	for _, r := range rs {
		if !r.Allowed {
			return r, true
		}
	}
	return RuleResult{}, false
}
func NormalizeStatus(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
func SortRecords(rs []Record) []Record {
	out := append([]Record(nil), rs...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func LatestRecord(rs []Record) (Record, bool) {
	if len(rs) == 0 {
		return Record{}, false
	}
	out := SortRecords(rs)
	return out[len(out)-1], true
}
func CountByStatus(rs []Record) map[string]int {
	out := map[string]int{}
	for _, r := range rs {
		out[r.Status]++
	}
	return out
}
func CountBySite(rs []Record) map[string]int {
	out := map[string]int{}
	for _, r := range rs {
		out[r.SiteID]++
	}
	return out
}
func CountByWorker(rs []Record) map[string]int {
	out := map[string]int{}
	for _, r := range rs {
		out[r.WorkerID]++
	}
	return out
}
func HasPending(rs []Record) bool {
	for _, r := range rs {
		if r.Status == StatusPending {
			return true
		}
	}
	return false
}
func HasSettled(rs []Record) bool {
	for _, r := range rs {
		if r.Status == StatusSettled {
			return true
		}
	}
	return false
}
func IsClosedAt(p Period, t time.Time) bool { return p.Closed && !p.ClosedAt.After(t) }
func PeriodDuration(p Period) time.Duration {
	if !p.Closed {
		return time.Since(p.ClosedAt)
	}
	return p.ClosedAt.Sub(time.Time{})
}
func MergePeriods(a, b Period) Period {
	if a.ID == "" {
		return ClonePeriodValue(b)
	}
	if b.ID == "" {
		return ClonePeriodValue(a)
	}
	a.RecordIDs = append(a.RecordIDs, b.RecordIDs...)
	a.Total += b.Total
	a.Closed = a.Closed || b.Closed
	return a
}
func ClonePeriodValue(p Period) Period { p.RecordIDs = append([]string(nil), p.RecordIDs...); return p }
func RemoveRecord(p Period, id string) Period {
	out := p.RecordIDs[:0]
	for _, v := range p.RecordIDs {
		if v != id {
			out = append(out, v)
		}
	}
	p.RecordIDs = out
	return p
}
func ContainsRecord(p Period, id string) bool {
	for _, v := range p.RecordIDs {
		if v == id {
			return true
		}
	}
	return false
}
func PeriodSize(p Period) int { return len(p.RecordIDs) }
func RecordIDs(rs []Record) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		out = append(out, r.ID)
	}
	return out
}
func StatusCounts(rs []Record) []string {
	m := CountByStatus(rs)
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
