package controllers

import (
	api "github.com/boardware-cloud/argus-api"
	argusServices "github.com/boardware-cloud/argus-service/services"
	coreServices "github.com/boardware-cloud/core/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

var accountService coreServices.AccountService
var argusService argusServices.ArgusService

func Init(inject *gorm.DB) {
	argusServices.Init(inject)
	coreServices.Init(inject)
	argusServices.Init(inject)
	accountService = coreServices.NewAccountService(inject)
	argusService = argusServices.NewArgusService(inject)
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
