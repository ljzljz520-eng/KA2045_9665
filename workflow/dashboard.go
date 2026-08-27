package workflow

import (
	"sort"
	"time"
	"workpay/domain"
)

type Dashboard struct {
	PeriodID  string
	State     string
	Total     float64
	Records   int
	UpdatedAt time.Time
}

func BuildDashboard(p domain.Period, rs []domain.Record) Dashboard {
	state := "open"
	if p.Closed {
		state = "closed"
	}
	return Dashboard{PeriodID: p.ID, State: state, Total: domain.SumRecords(rs), Records: len(rs), UpdatedAt: time.Now().UTC()}
}
func RefreshDashboard(old Dashboard, p domain.Period, rs []domain.Record) Dashboard {
	next := BuildDashboard(p, rs)
	if old.PeriodID == next.PeriodID && old.Records == next.Records && old.Total == next.Total {
		next.UpdatedAt = old.UpdatedAt
	}
	return next
}
func DashboardReady(d Dashboard) bool { return d.PeriodID != "" && d.Records >= 0 }
func DashboardLabel(d Dashboard) string {
	if d.State == "closed" {
		return "archived"
	}
	return "in-progress"
}
func RankWorkers(rs []domain.Record) []string {
	m := map[string]float64{}
	for _, r := range rs {
		m[r.WorkerID] += r.Amount()
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return m[keys[i]] > m[keys[j]] })
	return keys
}
func RankSites(rs []domain.Record) []string {
	m := map[string]float64{}
	for _, r := range rs {
		m[r.SiteID] += r.Amount()
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return m[keys[i]] > m[keys[j]] })
	return keys
}
func Window(rs []domain.Record, start, end time.Time) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if !r.CreatedAt.Before(start) && r.CreatedAt.Before(end) {
			out = append(out, r)
		}
	}
	return out
}
func Recent(rs []domain.Record, d time.Duration) []domain.Record {
	return Window(rs, time.Now().Add(-d), time.Now().Add(time.Second))
}
func PeriodReady(p domain.Period, rs []domain.Record) bool {
	if p.Closed || len(rs) != len(p.RecordIDs) {
		return false
	}
	for _, r := range rs {
		if !r.IsSettled() {
			return false
		}
	}
	return true
}
func Missing(p domain.Period, rs []domain.Record) []string {
	seen := map[string]bool{}
	for _, r := range rs {
		seen[r.ID] = true
	}
	out := []string{}
	for _, id := range p.RecordIDs {
		if !seen[id] {
			out = append(out, id)
		}
	}
	return out
}
