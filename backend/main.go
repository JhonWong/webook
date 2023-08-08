package main

import (
	"strings"
	"time"

	"github.com/JhonWong/webook/backend/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("this is first")
	})
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		//AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	u := web.NewUserHandler()
	u.RegisterRoutes(server)

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
