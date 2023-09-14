package async

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/johnwongx/webook/backend/internal/domain"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/service/sms/async/percent"

	"github.com/johnwongx/webook/backend/internal/service/sms"
)

type Service struct {
	svc        sms.Service
	threshod   float64
	pct        *percent.Percent[error]
	repo       repository.SMSRepository
	checkInter time.Duration

	resendCnt      int
	isCheckStarted int32
}

func NewService(svc sms.Service, threshod float64, errRange int, repo repository.SMSRepository,
	checkInter time.Duration) sms.Service {
	return &Service{
		svc:      svc,
		threshod: threshod,
		pct: percent.NewPercent[error](errRange, func(err error) bool {
			return err != nil
		}),
		repo:       repo,
		checkInter: checkInter,
		resendCnt:  1,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)

	s.pct.Add(err)
	if err != nil {
		//如果错误比例超过指定比例，认为第三方服务崩溃
		if s.isCrashed() {
			s.asyncStore(domain.SMSInfo{tpl, args, numbers})
		}
		return err
	}

	if !s.isCrashed() {
		//重新发送
		isStarted := atomic.LoadInt32(&s.isCheckStarted)
		if isStarted == 0 {
			if atomic.CompareAndSwapInt32(&s.isCheckStarted, isStarted, 1) {
				go s.checkService(ctx)
			}
		}

	}

	return nil
}

func (s *Service) asyncStore(info domain.SMSInfo) {
	go func() {
		//存储发送失败的信息
		s.repo.Put(info)
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
	if s.repo.IsEmpty() {
		return
	}

	infoArr, err := s.repo.Get(s.resendCnt)
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
			s.pct.Add(err)
			if err != nil {
				//再次发送失败
				s.asyncStore(info)
			}
		}(info, wg)
	}
	wg.Wait()

	if s.repo.IsEmpty() {
		s.resendCnt = 1
	} else {
		if !s.isCrashed() {
			if arrCnt == s.resendCnt {
				//说明失败的个数至少等于要发送的个数
				s.resendCnt *= 2
			}
		} else {
			s.resendCnt /= 2
		}
	}
}

func (s *Service) isCrashed() bool {
	return s.pct.Percent() > s.threshod
}
