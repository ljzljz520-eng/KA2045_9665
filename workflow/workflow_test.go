package workflow

import (
	"context"
	"path/filepath"
	"testing"
	"workpay/domain"
	"workpay/store"
)

func TestWorkflowOne(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	svc := New(s)
	if e := svc.Register(domain.NewProfile("w", "Worker", "worker", true), domain.NewRecord("r", "w", "site", 3, 10)); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowTwo(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	svc := New(s)
	svc.Register(domain.NewProfile("w", "Worker", "worker", true), domain.NewRecord("r", "w", "site", 3, 10))
	total, e := svc.Process(context.Background(), domain.Period{ID: "p", RecordIDs: []string{"r"}})
	if e != nil || total != 30 {
		t.Fatal(total, e)
	}
}
func TestWorkflowThree(t *testing.T) {
	if DescribePeriod(domain.Period{ID: "p"}) != "empty" {
		t.Fatal("state")
	}
}
