package controllers

import (
	"context"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine

var db *gorm.DB

func Init(inject context.Context) {
	db = inject.Value("db").(*gorm.DB)
	router = gin.Default()
	// router.Use(middleware.Auth())
	router.Use(middleware.CorsMiddleware())
	middleware.Health(router)
	var monitorApi = &MonitorApi{}
	api.MonitorApiInterfaceMounter(router, monitorApi)
	services.Init(inject)
}

func Run(addr ...string) {
	router.Run(addr...)
}
