package ioc

import (
	mysms "github.com/johnwongx/webook/backend/internal/service/sms"
	"github.com/johnwongx/webook/backend/internal/service/sms/localsms"
	"github.com/johnwongx/webook/backend/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitLocalSms() mysms.Service {
	return localsms.NewService()
}

func InitTencentSms() mysms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("SMS_SECRET_ID not found")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("SMS_SECRET_KEY not found")
	}

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}

	return tencent.NewService(c, "1400849905", "小葵花编程课堂个人公众号")
}
