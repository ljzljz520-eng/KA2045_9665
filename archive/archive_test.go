package archive

import (
	"path/filepath"
	"testing"
	"workpay/domain"
	"workpay/store"
)

func TestArchive(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := domain.NewRecord("r", "w", "s", 2, 4).MarkSettled()
	s.PutRecord(r)
	p := domain.Period{ID: "p", RecordIDs: []string{"r"}}
	got, e := New(s).ArchivePeriod(p, "admin")
	if e != nil || !got.Closed {
		t.Fatal(got, e)
	}
}
