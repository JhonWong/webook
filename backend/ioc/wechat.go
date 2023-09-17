package ioc

import (
	"net/http"
	"os"

	"github.com/johnwongx/webook/backend/internal/service/oauth2/wechat"
	"github.com/johnwongx/webook/backend/internal/web"
)

func InitWechatService() wechat.Service {
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("WECHAT_APP_ID not found")
	}
	secretKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("WECHAT_APP_SECRET not found")
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
		},
	}
	return wechat.NewService(appID, secretKey, client)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: true,
	}
}
