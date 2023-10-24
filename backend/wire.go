//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/repository/dao/article"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/internal/web/jwt"
	"github.com/johnwongx/webook/backend/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,

		dao.NewUserDAO,
		article.NewGORMArticleDAO,
		//article.NewGORMAuthorArticleDAO,
		//article.NewGORMReaderArticleDAO,

		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewArticleRepository,
		//repository.NewAuthorArticleRepository,
		//repository.NewReaderArticleRepository,

		ioc.InitTencentSms,
		ioc.InitWechatService,
		ioc.NewWechatHandlerConfig,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,

		web.NewUserHandler,
		web.NewWechatHandler,
		web.NewArticleHandler,
		jwt.NewRedisJwtHandler,

		ioc.InitRedisRateLimit,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return new(gin.Engine)
}
