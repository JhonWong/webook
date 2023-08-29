// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/JhonWong/webook/backend/internal/repository/cache"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/JhonWong/webook/backend/internal/service"
	"github.com/JhonWong/webook/backend/internal/web"
	"github.com/JhonWong/webook/backend/ioc"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	v := ioc.InitMiddlewares(cmdable)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	smsService := ioc.InitLocalSms()
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	codeService := service.NewCodeService(smsService, codeRepository)
	userHandler := web.NewUserHandler(userService, codeService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
