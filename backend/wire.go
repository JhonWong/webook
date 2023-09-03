//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/johnwongx/webook/backend/internal/repository"
	"github.com/johnwongx/webook/backend/internal/repository/cache"
	"github.com/johnwongx/webook/backend/internal/repository/dao"
	"github.com/johnwongx/webook/backend/internal/service"
	"github.com/johnwongx/webook/backend/internal/web"
	"github.com/johnwongx/webook/backend/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, ioc.InitLRUCache,

		dao.NewUserDAO,

		cache.NewUserCache,
		cache.NewLocalCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		ioc.InitLocalSms,

		service.NewUserService,
		service.NewCodeService,

		web.NewUserHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return new(gin.Engine)
}
