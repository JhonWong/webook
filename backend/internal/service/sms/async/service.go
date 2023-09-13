package async

import (
	"context"

	"github.com/johnwongx/webook/backend/internal/service/sms"
)

type Service struct {
	svc sms.Service
}

func NewService(svc sms.Service) sms.Service {
	return &Service{
		svc: svc,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string,
	numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)

	return nil
}
