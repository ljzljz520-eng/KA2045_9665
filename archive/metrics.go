package archive

import (
	"sort"
	"time"
	"workpay/domain"
)

type ArchiveMetrics struct {
	Records, Settled, Archived int
	Gross, ArchivedGross       float64
	Sites, Workers             int
}

func ComputeMetrics(rs []domain.Record) ArchiveMetrics {
	m := ArchiveMetrics{}
	sites := map[string]bool{}
	workers := map[string]bool{}
	for _, r := range rs {
		m.Records++
		m.Gross += r.Amount()
		sites[r.SiteID] = true
		workers[r.WorkerID] = true
		if r.Status == domain.StatusSettled {
			m.Settled++
		}
		if r.Status == domain.StatusArchived {
			m.Archived++
			m.ArchivedGross += r.Amount()
		}
	}
	m.Sites = len(sites)
	m.Workers = len(workers)
	return m
}
func Completion(m ArchiveMetrics) float64 {
	if m.Records == 0 {
		return 0
	}
	return float64(m.Archived) / float64(m.Records)
}
func SettlementRate(m ArchiveMetrics) float64 {
	if m.Records == 0 {
		return 0
	}
	return float64(m.Settled+m.Archived) / float64(m.Records)
}
func MetricsAt(rs []domain.Record, t time.Time) ArchiveMetrics {
	out := []domain.Record{}
	for _, r := range rs {
		if !r.CreatedAt.After(t) {
			out = append(out, r)
		}
	}
	return ComputeMetrics(out)
}
func StatusOrder() []string {
	return []string{domain.StatusPending, domain.StatusProcessing, domain.StatusSettled, domain.StatusArchived}
}
func SortByAmount(rs []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rs...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Amount() > out[j].Amount() })
	return out
}
func SortByWorker(rs []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rs...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].WorkerID < out[j].WorkerID })
	return out
}
func AmountRange(rs []domain.Record) (float64, float64) {
	if len(rs) == 0 {
		return 0, 0
	}
	min, max := rs[0].Amount(), rs[0].Amount()
	for _, r := range rs[1:] {
		if r.Amount() < min {
			min = r.Amount()
		}
		if r.Amount() > max {
			max = r.Amount()
		}
	}
	return min, max
}
func MedianAmount(rs []domain.Record) float64 {
	if len(rs) == 0 {
		return 0
	}
	v := SortByAmount(rs)
	n := len(v)
	if n%2 == 1 {
		return v[n/2].Amount()
	}
	return (v[n/2-1].Amount() + v[n/2].Amount()) / 2
}
func WorkerSummary(rs []domain.Record) map[string]ArchiveMetrics {
	out := map[string]ArchiveMetrics{}
	for _, r := range rs {
		m := out[r.WorkerID]
		m.Records++
		m.Gross += r.Amount()
		if r.Status == domain.StatusSettled {
			m.Settled++
		}
		if r.Status == domain.StatusArchived {
			m.Archived++
			m.ArchivedGross += r.Amount()
		}
		out[r.WorkerID] = m
	}
	return out
}
func SiteSummary(rs []domain.Record) map[string]ArchiveMetrics {
	out := map[string]ArchiveMetrics{}
	for _, r := range rs {
		m := out[r.SiteID]
		m.Records++
		m.Gross += r.Amount()
		if r.Status == domain.StatusArchived {
			m.Archived++
		}
		out[r.SiteID] = m
	}
	return out
}
func ActiveWorkers(ps []domain.Profile) []domain.Profile {
	out := []domain.Profile{}
	for _, p := range ps {
		if p.Active {
			out = append(out, p)
		}
	}
	return out
}
func Roles(ps []domain.Profile) map[string]int {
	out := map[string]int{}
	for _, p := range ps {
		out[p.Role]++
	}
	return out
}
func AuditTimeline(as []domain.Audit) []domain.Audit {
	out := append([]domain.Audit(nil), as...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func EventTimeline(es []domain.Event) []domain.Event {
	out := append([]domain.Event(nil), es...)
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out
}
func LatestAudit(as []domain.Audit) (domain.Audit, bool) {
	if len(as) == 0 {
		return domain.Audit{}, false
	}
	v := AuditTimeline(as)
	return v[len(v)-1], true
}
func LatestEvent(es []domain.Event) (domain.Event, bool) {
	if len(es) == 0 {
		return domain.Event{}, false
	}
	v := EventTimeline(es)
	return v[len(v)-1], true
}
func FilterAudits(as []domain.Audit, actor string) []domain.Audit {
	out := []domain.Audit{}
	for _, a := range as {
		if actor == "" || a.Actor == actor {
			out = append(out, a)
		}
	}
	return out
}
func FilterEvents(es []domain.Event, typ string) []domain.Event {
	out := []domain.Event{}
	for _, e := range es {
		if typ == "" || e.Type == typ {
			out = append(out, e)
		}
	}
	return out
}
