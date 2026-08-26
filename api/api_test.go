package api

import (
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"workpay/store"
	"workpay/workflow"
)

func TestHTTPRegister(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	req := httptest.NewRequest("POST", "/records", strings.NewReader(`{"Profile":{"ID":"w","Name":"W","Active":true},"Record":{"ID":"r","WorkerID":"w","SiteID":"s","Quantity":1,"Rate":2}}`))
	w := httptest.NewRecorder()
	New(workflow.New(s)).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
