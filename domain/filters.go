package domain

type Filter struct {
	WorkerID, SiteID, Status string
	MinAmount                float64
}

func (f Filter) Match(r Record) bool {
	if f.WorkerID != "" && f.WorkerID != r.WorkerID {
		return false
	}
	if f.SiteID != "" && f.SiteID != r.SiteID {
		return false
	}
	if f.Status != "" && f.Status != r.Status {
		return false
	}
	return r.Amount() >= f.MinAmount
}
func FilterRecords(rs []Record, f Filter) []Record {
	out := []Record{}
	for _, r := range rs {
		if f.Match(r) {
			out = append(out, r)
		}
	}
	return out
}
func GroupBySite(rs []Record) map[string][]Record {
	out := map[string][]Record{}
	for _, r := range rs {
		out[r.SiteID] = append(out[r.SiteID], r)
	}
	return out
}
func GroupByWorker(rs []Record) map[string][]Record {
	out := map[string][]Record{}
	for _, r := range rs {
		out[r.WorkerID] = append(out[r.WorkerID], r)
	}
	return out
}
