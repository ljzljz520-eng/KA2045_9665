package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"workpay/domain"
	"workpay/store"
	"workpay/workflow"
)

func TestBusinessChain18(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	svc := workflow.New(s)
	if e := svc.Register(domain.NewProfile("w", "Worker", "worker", true), domain.NewRecord("r", "w", "site", 2, 7)); e != nil {
		t.Fatal(e)
	}
	p, _ := svc.Registry.CreatePeriod("p", []string{"r"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := svc.CancelAndClose(ctx, p, "admin"); e == nil {
		t.Fatal("cancel should stop processing")
	}
	time.Sleep(10 * time.Millisecond)
	if svc.Calculator.ResourceCount() != 0 {
		t.Fatalf("calculator resources remain occupied: %d", svc.Calculator.ResourceCount())
	}
}
