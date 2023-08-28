package localsms

import (
	"context"
	"fmt"
)

type LocalService struct {
}

func NewService() *LocalService {
	return &LocalService{}
}

func (s *LocalService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	fmt.Println(args)
	return nil
}
