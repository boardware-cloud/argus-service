package controllers

import (
	"time"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/common/server"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/core/controllers"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

var middleware controllers.Middleware

func timeoutTesting(c *gin.Context) {
	duration := time.Duration(utils.StringToUint(c.Param("duration")) * uint(time.Second))
	time.Sleep(duration)
}

func init() {
	router = gin.Default()
	router.Use(server.CorsMiddleware())
	api.MonitorApiInterfaceMounter(router, &MonitorApi{})
	router.GET("/testing/delay/:duration", timeoutTesting)
	router.GET("/testing/down", func(ctx *gin.Context) {
	})
}

func Run(addr ...string) {
	router.Run(addr...)
}
