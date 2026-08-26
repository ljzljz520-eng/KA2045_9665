package store

import (
	"path/filepath"
	"testing"
	"workpay/domain"
)

func TestStoreRoundTrip(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r := domain.NewRecord("r", "w", "s", 1, 2)
	if e = s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	got, e := s.GetRecord("r")
	if e != nil || got.ID != "r" {
		t.Fatal(got, e)
	}
}
