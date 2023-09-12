package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/johnwongx/webook/backend/internal/service/sms"
)

type Service struct {
	svc sms.Service
	key string
}

func NewService(svc sms.Service, key string) *Service {
	return &Service{
		svc: svc,
		key: key,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string,
	numbers ...string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tpl, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return err
	}

	if token == nil || !token.Valid {
		return errors.New("token非法")
	}

	return s.svc.Send(ctx, claims.tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	tpl string
}
