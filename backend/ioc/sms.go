package ioc

import (
	mysms "github.com/johnwongx/webook/backend/internal/service/sms"
	"github.com/johnwongx/webook/backend/internal/service/sms/localsms"
	smslimit "github.com/johnwongx/webook/backend/internal/service/sms/ratelimit"
	"github.com/johnwongx/webook/backend/internal/service/sms/tencent"
	"github.com/johnwongx/webook/backend/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"time"
)

func InitLocalSms() mysms.Service {
	return localsms.NewService()
}

func InitTencentSms(redisClient redis.Cmdable) mysms.Service {
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
	svc := tencent.NewService(c, "1400849905", "小葵花编程课堂个人公众号")
	limiter := ratelimit.NewRedisSliderWindowLimiter(redisClient, time.Second, 3000)
	return smslimit.NewServiceSMSRateLimiter(svc, limiter)
}
