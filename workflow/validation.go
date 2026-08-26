package workflow

import (
	"fmt"
	"workpay/domain"
)

func ValidateRegistration(p domain.Profile, r domain.Record) error {
	if e := domain.ValidateProfile(p); e != nil {
		return e
	}
	if e := r.Validate(p); e != nil {
		return e
	}
	return nil
}
func ValidateProcessing(p domain.Period) error {
	if !domain.CanProcess(p) {
		return fmt.Errorf("period cannot process")
	}
	return nil
}
func ValidateClosing(p domain.Period, rs []domain.Record) error { return domain.ValidateArchive(p, rs) }
func ExplainError(e error) string {
	if e == nil {
		return ""
	}
	if e == domain.ErrPeriodClosed {
		return "period already closed"
	}
	return e.Error()
}
func StatusForError(e error) int {
	if e == nil {
		return 200
	}
	if e == domain.ErrPeriodClosed {
		return 409
	}
	return 400
}
func EnsureWorkflowInputs(p domain.Profile, r domain.Record, period domain.Period) error {
	if e := ValidateRegistration(p, r); e != nil {
		return e
	}
	if period.ID == "" {
		return fmt.Errorf("period id required")
	}
	return nil
}
func RecordInPeriod(p domain.Period, r domain.Record) bool { return domain.ContainsRecord(p, r.ID) }
func CompletionMessage(p domain.Period) string {
	if p.Closed {
		return fmt.Sprintf("period %s complete", p.ID)
	}
	return fmt.Sprintf("period %s pending", p.ID)
}
func IsRetryable(e error) bool           { return e != nil && e != domain.ErrPeriodClosed }
func WorkflowReady(p domain.Period) bool { return domain.CanProcess(p) }
