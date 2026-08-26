package payroll

import (
	"context"
	"path/filepath"
	"testing"
	"workpay/domain"
	"workpay/store"
)

func TestSettle(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	s.PutRecord(domain.NewRecord("r", "w", "s", 2, 4))
	c := New(s, 2)
	total, e := c.SettlePeriod(context.Background(), domain.Period{ID: "p", RecordIDs: []string{"r"}})
	if e != nil || total != 8 {
		t.Fatal(total, e)
	}
}
func TestPolicy(t *testing.T) {
	r := domain.NewRecord("r", "w", "s", 10, 2)
	if DefaultPolicy().Amount(r) <= r.Amount() {
		t.Fatal("overtime")
	}
}
