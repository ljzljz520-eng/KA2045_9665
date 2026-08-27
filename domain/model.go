package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidRecord   = errors.New("invalid piecework record")
	ErrInactiveProfile = errors.New("inactive worker profile")
	ErrPeriodClosed    = errors.New("period is closed")
)

type Record struct {
	ID, WorkerID, SiteID string
	Quantity, Rate       float64
	Status               string
	CreatedAt            time.Time
}
type Profile struct {
	ID, Name, Role string
	Active         bool
}
type Event struct {
	ID, RecordID, Type, Payload string
	CreatedAt                   time.Time
}
type Audit struct {
	ID, Actor, Action, Target string
	CreatedAt                 time.Time
}
type Period struct {
	ID        string
	RecordIDs []string
	Closed    bool
	Total     float64
	ClosedAt  time.Time
}

func (r Record) Validate(p Profile) error {
	if r.ID == "" || r.WorkerID == "" || r.SiteID == "" {
		return ErrInvalidRecord
	}
	if r.Quantity <= 0 || r.Rate <= 0 {
		return ErrInvalidRecord
	}
	if !p.Active {
		return ErrInactiveProfile
	}
	return nil
}
func (r Record) Amount() float64      { return r.Quantity * r.Rate }
func (r Record) IsSettled() bool      { return r.Status == "settled" || r.Status == "archived" }
func (r Record) MarkSettled() Record  { r.Status = "settled"; return r }
func (r Record) MarkArchived() Record { r.Status = "archived"; return r }
func NewRecord(id, worker, site string, qty, rate float64) Record {
	return Record{ID: id, WorkerID: worker, SiteID: site, Quantity: qty, Rate: rate, Status: "pending", CreatedAt: time.Now().UTC()}
}
func NewProfile(id, name, role string, active bool) Profile {
	return Profile{ID: id, Name: name, Role: role, Active: active}
}
func NewEvent(id, rid, typ, payload string) Event {
	return Event{ID: id, RecordID: rid, Type: typ, Payload: payload, CreatedAt: time.Now().UTC()}
}
func NewAudit(id, actor, action, target string) Audit {
	return Audit{ID: id, Actor: actor, Action: action, Target: target, CreatedAt: time.Now().UTC()}
}
func ValidatePeriod(p Period) error {
	if p.ID == "" {
		return fmt.Errorf("period id required")
	}
	if p.Closed {
		return ErrPeriodClosed
	}
	return nil
}
func (p Period) Add(id string) Period { p.RecordIDs = append(p.RecordIDs, id); return p }
func (p Period) Close(total float64) Period {
	p.Closed = true
	p.Total = total
	p.ClosedAt = time.Now().UTC()
	return p
}
