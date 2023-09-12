package tencent

import (
	"context"
	"fmt"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentService struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId, signName string) *TencentService {
	return &TencentService{
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		client:   client,
	}
}

func (s *TencentService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(number)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Message) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *TencentService) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
