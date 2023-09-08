package ratelimit

import (
	"context"
	"fmt"
	"github.com/johnwongx/webook/backend/internal/service/sms"
	"github.com/johnwongx/webook/backend/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发限流")

type ServiceSMSRateLimiter struct {
	Svc     sms.Service
	Limiter ratelimit.Limiter
}

func NewServiceSMSRateLimiter(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &ServiceSMSRateLimiter{
		Svc:     svc,
		Limiter: limiter,
	}
}

func (s *ServiceSMSRateLimiter) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.Limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		//限流
		//保守策略：如果下游比较脆弱
		//开放策略：下游比较稳定，业务可用性要求高, 尽量容错
		return fmt.Errorf("短信服务判断出错，%w", err)
	}
	if limited {
		return errLimited
	}

	err = s.Svc.Send(ctx, tpl, args, numbers...)
	return err
}
