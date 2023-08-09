package main

import (
	"strings"
	"time"

	"github.com/JhonWong/webook/backend/internal/repository"
	"github.com/JhonWong/webook/backend/internal/repository/dao"
	"github.com/JhonWong/webook/backend/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/JhonWong/webook/backend/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("this is first")
	})

	//跨域问题
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	dsn := "root:root@tcp(localhost:13316)/webook"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	dao := dao.NewUserDAO(db)
	r := repository.NewUserRepository(dao)
	us := service.NewUserService(r)
	u := web.NewUserHandler(us)
	u.RegisterRoutes(server)

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
