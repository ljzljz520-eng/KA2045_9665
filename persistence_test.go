package main

import (
	"path/filepath"
	"testing"
	"workpay/domain"
	"workpay/store"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "db")
	s, e := store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	r := domain.NewRecord("r", "w", "s", 1, 3)
	if e = s.PutRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetRecord("r")
	if e != nil || got.ID != "r" {
		t.Fatal(got, e)
	}
}
