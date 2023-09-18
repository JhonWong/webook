//go:build wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,

		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		ioc.InitTencentSms,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,

		service.NewUserService,
		service.NewCodeService,

		web.NewUserHandler,
		web.NewWechatHandler,
		jwt.NewRedisJwtHandler,

		ioc.InitRedisRateLimit,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return new(gin.Engine)
}
