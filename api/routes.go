package api

import "net/http"

func Routes(h *Handler) *http.ServeMux { m := http.NewServeMux(); m.Handle("/", h); return m }
