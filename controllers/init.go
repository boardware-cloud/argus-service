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
	var monitorApi = &MonitorApi{}
	api.MonitorApiInterfaceMounter(router, monitorApi)
	var reservedApi = &ReservedApi{}
	api.ReservedApiInterfaceMounter(router, reservedApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
