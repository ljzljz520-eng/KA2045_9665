package domain

import "time"

type Change struct {
	Field, Before, After string
	At                   time.Time
}

func NewChange(field, before, after string) Change {
	return Change{field, before, after, time.Now().UTC()}
}
func ChangesForRecord(before, after Record) []Change {
	out := []Change{}
	if before.Quantity != after.Quantity {
		out = append(out, NewChange("quantity", formatFloat(before.Quantity), formatFloat(after.Quantity)))
	}
	if before.Rate != after.Rate {
		out = append(out, NewChange("rate", formatFloat(before.Rate), formatFloat(after.Rate)))
	}
	if before.Status != after.Status {
		out = append(out, NewChange("status", before.Status, after.Status))
	}
	return out
}
func formatFloat(v float64) string {
	if v == 0 {
		return "0"
	}
	return time.Duration(v).String()
}
func Changed(before, after Record) bool    { return len(ChangesForRecord(before, after)) > 0 }
func ChangeCount(before, after Record) int { return len(ChangesForRecord(before, after)) }
func EventForChange(id string, c Change) Event {
	return NewEvent(id, "", "changed", c.Field+":"+c.Before+">"+c.After)
}
func AuditForChange(id, actor, target string, c Change) Audit {
	return NewAudit(id, actor, "change:"+c.Field, target)
}
func WithinWindow(c Change, start, end time.Time) bool {
	return !c.At.Before(start) && c.At.Before(end)
}
func LatestChange(cs []Change) (Change, bool) {
	if len(cs) == 0 {
		return Change{}, false
	}
	latest := cs[0]
	for _, c := range cs[1:] {
		if c.At.After(latest.At) {
			latest = c
		}
	}
	return latest, true
}
func ChangesSince(cs []Change, t time.Time) []Change {
	out := []Change{}
	for _, c := range cs {
		if !c.At.Before(t) {
			out = append(out, c)
		}
	}
	return out
}
