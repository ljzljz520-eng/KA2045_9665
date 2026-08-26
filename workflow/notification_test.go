package workflow

import "testing"

func TestNotifier(t *testing.T) {
	n := NewNotifier()
	if e := n.Send("admin", "ready"); e != nil || len(n.Sent()) != 1 {
		t.Fatal(e)
	}
}
