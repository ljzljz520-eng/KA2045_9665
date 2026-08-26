package registry

import (
	"path/filepath"
	"testing"
	"workpay/domain"
	"workpay/store"
)

func TestRegister(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := New(s)
	if e := r.RegisterProfile(domain.NewProfile("w", "W", "worker", true)); e != nil {
		t.Fatal(e)
	}
	if e := r.RegisterRecord(domain.NewRecord("r", "w", "s", 4, 5)); e != nil {
		t.Fatal(e)
	}
}
