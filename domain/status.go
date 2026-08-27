package domain

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusSettled    = "settled"
	StatusArchived   = "archived"
)

func StatusAllowed(s string) bool {
	switch s {
	case StatusPending, StatusProcessing, StatusSettled, StatusArchived:
		return true
	default:
		return false
	}
}
func Transition(from, to string) bool {
	if !StatusAllowed(from) || !StatusAllowed(to) {
		return false
	}
	if from == to {
		return true
	}
	return (from == StatusPending && to == StatusProcessing) || (from == StatusProcessing && to == StatusSettled) || (from == StatusSettled && to == StatusArchived)
}
