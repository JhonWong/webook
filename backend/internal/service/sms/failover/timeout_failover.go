package failover

import (
	"context"
	"sync/atomic"

	"github.com/johnwongx/webook/backend/internal/service/sms"
)

type TimeoutFailover struct {
	svcs []sms.Service
	//当前服务编号
	idx uint32
	//连续失败次数
	cnt uint32
	//连续失败次数阈值
	threshold uint32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold uint32) *TimeoutFailover {
	return &TimeoutFailover{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (s *TimeoutFailover) Send(ctx context.Context, tpl string,
	args []string, numbers ...string) error {
	//检查是否超过阈值，超过时切换下一个
	idx := atomic.LoadUint32(&s.idx)
	cnt := atomic.LoadUint32(&s.cnt)
	if cnt > s.threshold {
		newIdx := (idx + 1) % uint32(len(s.svcs))
		if atomic.CompareAndSwapUint32(&s.idx, idx, newIdx) {
			atomic.StoreUint32(&s.cnt, 0)
		}

		idx = atomic.LoadUint32(&s.idx)
	}

	svc := s.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case nil:
		//中断连续计数
		atomic.StoreUint32(&s.cnt, uint32(0))
		return nil
	case context.DeadlineExceeded:
		//增加连续计数
		atomic.AddUint32(&s.cnt, uint32(1))
		return err
	default:
		//其他错误
		return err
	}
}
