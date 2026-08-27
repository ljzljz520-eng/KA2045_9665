package workflow

import (
	"context"
	"workpay/archive"
	"workpay/domain"
	"workpay/payroll"
	"workpay/registry"
	"workpay/store"
)

type Service struct {
	Registry   *registry.Registry
	Calculator *payroll.Calculator
	Archive    *archive.Service
}

func New(db *store.BoltStore) *Service {
	return &Service{Registry: registry.New(db), Calculator: payroll.New(db, 4), Archive: archive.New(db)}
}
func (s *Service) Register(p domain.Profile, r domain.Record) error {
	if e := s.Registry.RegisterProfile(p); e != nil {
		return e
	}
	return s.Registry.RegisterRecord(r)
}
func (s *Service) Process(ctx context.Context, p domain.Period) (float64, error) {
	return s.Calculator.SettlePeriod(ctx, p)
}
func (s *Service) Close(p domain.Period, actor string) (domain.Period, error) {
	return s.Archive.ArchivePeriod(p, actor)
}
func (s *Service) CancelAndClose(ctx context.Context, p domain.Period, actor string) error {
	_, e := s.Process(ctx, p)
	if e != nil {
		return e
	}
	_, e = s.Close(p, actor)
	return e
}
