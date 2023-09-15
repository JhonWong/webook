package async

import (
	"context"
	"github.com/johnwongx/webook/backend/internal/service/sms/async/serviceprobe"
	"sync"
	"sync/atomic"
	"time"

	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/service/sms"
)

type Service struct {
	svc        sms.Service
	svcProbe   serviceprobe.ServiceProbe
	repo       repository.SMSRepository
	checkInter time.Duration

	resendCnt      int
	isCheckStarted int32
}

func NewService(svc sms.Service, svcProbe serviceprobe.ServiceProbe, repo repository.SMSRepository,
	checkInter time.Duration) sms.Service {
	return &Service{
		svc:        svc,
		svcProbe:   svcProbe,
		repo:       repo,
		checkInter: checkInter,
		resendCnt:  1,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)

	s.svcProbe.Add(ctx, err)
	if err != nil {
		if s.svcProbe.IsCrashed(ctx) {
			s.asyncStore(ctx, domain.SMSInfo{tpl, args, numbers})
		}
		return err
	}

	if !s.svcProbe.IsCrashed(ctx) {
		//重新发送
		isStarted := atomic.LoadInt32(&s.isCheckStarted)
		//重发携程未启动，并且存储的有错误信息
		if isStarted == 0 && !s.repo.IsEmpty(ctx) {
			if atomic.CompareAndSwapInt32(&s.isCheckStarted, isStarted, 1) {
				go s.checkService(ctx)
			}
		}

	}

	return nil
}

func (s *Service) asyncStore(ctx context.Context, info domain.SMSInfo) {
	go func() {
		//存储发送失败的信息
		err := s.repo.Put(ctx, info)
		if err != nil {
			//TODO:add log
		}
	}()
}

func (s *Service) checkService(ctx context.Context) {
	ticker := time.NewTicker(s.checkInter)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.asyncSend(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) asyncSend(ctx context.Context) {
	if s.repo.IsEmpty(ctx) {
		return
	}

	infoArr, err := s.repo.Get(ctx, s.resendCnt)
	if err != nil {
		return
	}

	arrCnt := len(infoArr)
	var wg sync.WaitGroup
	for _, info := range infoArr {
		wg.Add(1)
		go func(info domain.SMSInfo, wg sync.WaitGroup) {
			defer wg.Done()

			err = s.svc.Send(ctx, info.Tpl, info.Args, info.Numbers...)
			s.svcProbe.Add(ctx, err)
			if err != nil {
				//再次发送失败
				s.asyncStore(ctx, info)
			}
		}(info, wg)
	}
	wg.Wait()

	if s.repo.IsEmpty(ctx) {
		s.resendCnt = 1
	} else {
		if !s.svcProbe.IsCrashed(ctx) {
			if arrCnt == s.resendCnt {
				s.resendCnt *= 2
			}
		} else {
			s.resendCnt /= 2
		}
	}
}
