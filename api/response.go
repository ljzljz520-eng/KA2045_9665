package api

import (
	"net/http"
	"strconv"
	"time"
)

type ResponseMeta struct {
	RequestID string
	StartedAt time.Time
	Duration  time.Duration
}

func NewMeta(id string) ResponseMeta        { return ResponseMeta{RequestID: id, StartedAt: time.Now().UTC()} }
func (m ResponseMeta) Finish() ResponseMeta { m.Duration = time.Since(m.StartedAt); return m }
func (m ResponseMeta) Headers() map[string]string {
	return map[string]string{"X-Request-ID": m.RequestID, "X-Elapsed": m.Duration.String()}
}
func ApplyMeta(w http.ResponseWriter, m ResponseMeta) {
	for k, v := range m.Headers() {
		w.Header().Set(k, v)
	}
}
func ParseLimit(v string, defaultValue, max int) int {
	n, e := strconv.Atoi(v)
	if e != nil || n <= 0 {
		return defaultValue
	}
	if n > max {
		return max
	}
	return n
}
func ParseOffset(v string) int {
	n, e := strconv.Atoi(v)
	if e != nil || n < 0 {
		return 0
	}
	return n
}
func Paginate[T any](items []T, offset, limit int) []T {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = len(items)
	}
	if offset >= len(items) {
		return []T{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}
func Clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func BoolParam(v string) bool { return v == "1" || v == "true" || v == "yes" }
func ContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
