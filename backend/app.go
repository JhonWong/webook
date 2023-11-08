package main

import (
	"github.com/gin-gonic/gin"
	"github.com/johnwongx/webook/backend/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
