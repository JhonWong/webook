package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/JhonWong/webook/backend/internal/service/sms"
)

const codeTplId = "1907519"

type CodeService struct {
	svc        sms.Service
	repo       *repository.CodeRepository
	expiration time.Duration
}

func NewCodeService(svc sms.Service, repo *repository.CodeRepository) *CodeService {
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
	err := s.repo.Store(ctx, biz, phone, code, s.expiration)
	if err != nil {
		return err
	}

	//3.发送验证码
	expVal := fmt.Sprintf("%d", s.expiration.Minutes())
	params := []string{code, expVal}
	err = s.svc.Send(ctx, codeTplId, params, phone)
	return err
}

func (s *CodeService) Verify(ctx context.Context, biz, code, phone string) (bool, error) {
	return s.repo.Verify(ctx, biz, phone, code)
}

func (s *CodeService) generateCode(biz, phone string) string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%d", num)
}
