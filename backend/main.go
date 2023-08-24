package main

import (
	"strings"
	"time"

	"github.com/JhonWong/webook/backend/config"

	"github.com/JhonWong/webook/backend/pkg/ginx/middlewares/ratelimit"

	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/JhonWong/webook/backend/internal/repository/cache"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/JhonWong/webook/backend/internal/service"
	"github.com/JhonWong/webook/backend/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/JhonWong/webook/backend/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	db := initDB()
	rdb := initRedis()
	server := initServer()

	u := initUser(db, rdb)
	u.RegisterRoutes(server)

	server.Run(":8081")
}

func initServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("this is first")
	})

	//跨域问题
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	//store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	//	[]byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	//server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePath("/users/signup").
		IgnorePath("/users/login").Builder())
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePath("/users/signup").
	//	IgnorePath("/users/login").Builder())

	return server
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	dao := dao.NewUserDAO(db)
	cache := cache.NewUserCache(rdb)
	r := repository.NewUserRepository(dao, cache)
	us := service.NewUserService(r)
	u := web.NewUserHandler(us)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}
