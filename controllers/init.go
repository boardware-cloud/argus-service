package controllers

import (
	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/common/server"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	middleware.Health(router)
	api.MonitorApiInterfaceMounter(router, &MonitorApi{})
}

func Run(addr ...string) {
	router.Run(addr...)
}
