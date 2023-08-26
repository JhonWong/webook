package service

import (
	"context"
	"errors"
	"github.com/JhonWong/webook/backend/internal/repository"
	mysms "github.com/JhonWong/webook/backend/internal/service/sms"
	"time"
)

var (
	ErrSendTooMuch = repository.ErrSendTooMuch
)

type CodeService struct {
	svc        *mysms.Service
	repo       *repository.CodeRepository
	expiration time.Duration
}

func NewCodeService(svc *mysms.Service, repo *repository.CodeRepository) *CodeService {
	return &CodeService{
		svc:        svc,
		repo:       repo,
		expiration: time.Minute * 30,
	}
}

func (s *CodeService) Send(ctx context.Context, biz, phone string) error {
	//1.生成验证码
	code := s.generateCode(biz, phone)

	//2.存储验证码
	err := s.repo.Set(ctx, biz, code, phone)
	if err == repository.ErrSendTooMuch {
		return ErrSendTooMuch
	} else if err != nil {
		return errors.New("系统错误")
	}

	//3.发送验证码
	tplId := "1907519"
	//TODO:将过期时间转化为数字
	params := []string{code, "233"}
	err = s.svc.Send(ctx, tplId, params, params)
	return err
}

func (s *CodeService) Verify(ctx context.Context, biz, code, phone string) error {
	realCode, err := s.repo.Get(ctx, biz, phone)
	if err != nil {
		return err
	}

	if realCode != code {
		return err
	}
	return nil
}

func (s *CodeService) generateCode(biz, phone string) string {

}
