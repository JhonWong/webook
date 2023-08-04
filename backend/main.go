package main

import (
	"github.com/JhonWong/webook/backend/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	u := web.NewUserHandler()
	u.RegisterRoutes(server)

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
