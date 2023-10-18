package async2

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/service/sms"
	"github.com/johnwongx/webook/backend/pkg/logger"
	"time"
)

var _ sms.Service = &AsyncService{}

type AsyncService struct {
	svc  sms.Service
	repo repository.SMSAsyncRepository
	l    logger.Logger
}

func NewAsyncService(svc sms.Service, repo repository.SMSAsyncRepository, l logger.Logger) *AsyncService {
	return &AsyncService{
		svc:  svc,
		repo: repo,
		l:    l,
	}
}

func (s *AsyncService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	if s.NeedAsync() {
		// 将数据存储
		err := s.repo.Add(ctx, domain.SMSAsyncInfo{
			Tpl:           tpl,
			Args:          args,
			Numbers:       numbers,
			MaxRetryCount: 3,
		})
		return err
	}

	err := s.svc.Send(ctx, tpl, args, numbers...)
	s.UpdateSvcStatus(err)
	return err
}

// 主线程退出时，该异步退出
func (s *AsyncService) StartAsync() {
	for {
		s.SendAsync()
	}
}

func (s *AsyncService) SendAsync() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 抢占短信
	info, err := s.repo.PreemptWaitingSMS(ctx)
	cancel()
	switch err {
	case nil:
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = s.svc.Send(ctx, info.Tpl, info.Args, info.Numbers...)
		if err != nil {
			// 不做任何处理，因为取出数据时已经将短信的重试次数做了修改
			s.l.Error("Send sms info failed!",
				logger.Error(err),
				logger.Int64("id", info.Id))
		}

		res := err == nil
		err = s.repo.ReportScheduleResult(ctx, info.Id, res)
		if err != nil {
			s.l.Error("Update sms info status failed!",
				logger.Error(err))
		}
	case repository.ErrorEmptySmsInfo:
		// 未获取任何重发短信，等待一秒
		time.Sleep(time.Second)
	default:
		// 发生未知错误
		time.Sleep(time.Second)
		s.l.Error("GetWaitingSMS async sms info failed!",
			logger.Error(err))
	}
}

func (s *AsyncService) NeedAsync() bool {
	// TODO:
	return false
}

func (s *AsyncService) UpdateSvcStatus(err error) {

}
