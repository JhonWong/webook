package startup

import (
	"github.com/johnwongx/webook/backend/internal/service/oauth2/wechat"
	"github.com/johnwongx/webook/backend/pkg/logger"
)

func InitPhantomWechatService(l logger.Logger) wechat.Service {
	return wechat.NewService("", "", nil, l)
}
