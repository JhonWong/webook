package repository

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/sqlx"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
)

var _ SMSAsyncRepository = &smsAsyncRepository{}

var (
	ErrorEmptySmsInfo = errors.New("No sms info")
)

type SMSAsyncRepository interface {
	Add(ctx context.Context, info domain.SMSAsyncInfo) error
	PreemptWaitingSMS(ctx context.Context) (domain.SMSAsyncInfo, error)
	ReportScheduleResult(ctx context.Context, id int64, res bool) error
}

type smsAsyncRepository struct {
	d dao.AsyncSMSDAO
}

func NewSmsAsyncRepository(d dao.AsyncSMSDAO) *smsAsyncRepository {
	return &smsAsyncRepository{
		d: d,
	}
}

func (s *smsAsyncRepository) Add(ctx context.Context, info domain.SMSAsyncInfo) error {
	return s.d.Insert(ctx, dao.SMSAsyncInfo{
		Config: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				Tpl:     info.Tpl,
				Args:    info.Args,
				Numbers: info.Numbers,
			},
			Valid: true,
		},
		RetryMax: info.MaxRetryCount,
	})
}

func (s *smsAsyncRepository) PreemptWaitingSMS(ctx context.Context) (domain.SMSAsyncInfo, error) {
	info, err := s.d.GetWaitingSMS(ctx)
	if err != nil {
		return domain.SMSAsyncInfo{}, err
	}
	return domain.SMSAsyncInfo{
		Id:      info.Id,
		Tpl:     info.Config.Val.Tpl,
		Args:    info.Config.Val.Args,
		Numbers: info.Config.Val.Numbers,
	}, nil
}

func (s *smsAsyncRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return s.d.MarkSuccess(ctx, id)
	}
	return s.d.MarkFailed(ctx, id)
}
