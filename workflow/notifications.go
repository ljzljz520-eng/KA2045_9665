package workflow

import (
	"fmt"
	"sync"
	"workpay/domain"
)

type Notification struct{ Recipient, Message string }
type Notifier struct {
	mu   sync.Mutex
	sent []Notification
}

func NewNotifier() *Notifier { return &Notifier{sent: []Notification{}} }
func (n *Notifier) Send(recipient, message string) error {
	if recipient == "" || message == "" {
		return fmt.Errorf("notification fields required")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.sent = append(n.sent, Notification{recipient, message})
	return nil
}
func (n *Notifier) Sent() []Notification {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]Notification(nil), n.sent...)
}
func NotifyPeriod(n *Notifier, p domain.Period, recipient string) error {
	state := "open"
	if p.Closed {
		state = "closed"
	}
	return n.Send(recipient, fmt.Sprintf("period %s is %s", p.ID, state))
}
