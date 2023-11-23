package controllers

import (
	api "github.com/boardware-cloud/argus-api"
	argusServices "github.com/boardware-cloud/argus-service/services"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

var accountService = coreServices.GetAccountService()
var argusService = argusServices.GetArgusService()

func init() {
	router = gin.Default()
	router.Use(accountService.Auth())
	router.Use(middleware.CorsMiddleware())
	middleware.Health(router)
	var monitorApi = &MonitorApi{}
	api.MonitorApiInterfaceMounter(router, monitorApi)
}

func Run(addr ...string) {
	router.Run(addr...)
}
