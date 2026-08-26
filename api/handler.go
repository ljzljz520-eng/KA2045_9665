package api

import (
	"context"
	"encoding/json"
	"net/http"
	"workpay/domain"
	"workpay/workflow"
)

type Handler struct{ Service *workflow.Service }

func New(s *workflow.Service) *Handler { return &Handler{Service: s} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/health":
		h.health(w)
	case "/records":
		h.records(w, r)
	case "/periods/process":
		h.process(w, r)
	default:
		http.NotFound(w, r)
	}
}
func (h *Handler) health(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
func (h *Handler) records(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method", 405)
		return
	}
	var in struct {
		Profile domain.Profile
		Record  domain.Record
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		http.Error(w, "json", 400)
		return
	}
	if e := h.Service.Register(in.Profile, in.Record); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	writeJSON(w, in.Record)
}
func (h *Handler) process(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method", 405)
		return
	}
	var p domain.Period
	if json.NewDecoder(r.Body).Decode(&p) != nil {
		http.Error(w, "json", 400)
		return
	}
	total, e := h.Service.Process(r.Context(), p)
	if e != nil {
		http.Error(w, e.Error(), 409)
		return
	}
	writeJSON(w, map[string]any{"period": p.ID, "total": total})
}
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func RequestContext(r *http.Request) context.Context { return r.Context() }
