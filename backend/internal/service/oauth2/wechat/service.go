package wechat

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/domain"
	uuid "github.com/lithammer/shortuuid/v4"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx *gin.Context) (string, error)
}

type service struct {
	appID string
}

func NewWechatService(appID string) Service {
	return &service{
		appID: appID,
	}
}

func (s *service) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=SCOPE&state=STATE#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, s.appID, redirectURI, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {

}
