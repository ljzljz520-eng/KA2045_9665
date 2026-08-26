package domain

import "testing"

func TestRecordValidation(t *testing.T) {
	r := NewRecord("r1", "w1", "s1", 2, 3)
	if e := r.Validate(NewProfile("w1", "A", "worker", true)); e != nil {
		t.Fatal(e)
	}
	if r.Amount() != 6 {
		t.Fatal(r.Amount())
	}
}
func TestCalendar(t *testing.T) {
	if NewCalendar(2026, 2).Days != 28 {
		t.Fatal("days")
	}
}
