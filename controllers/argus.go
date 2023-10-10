package controllers

import (
	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
)

const DEFAULT_TIMEOUT = 10
const DEFAULT_INTERVAL = 5 * 60

type MonitorApi struct{}

func (MonitorApi) UpdateMonitor(c *gin.Context, monitorId string, putMonitorRequest api.PutMonitorRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			// services.UpdateMonitor(
			// 	account.ID,
			// 	utils.StringToUint(monitorId),
			// 	PutMonitorForward(putMonitorRequest),
			// ).Just(func(data services.Monitor) {
			// 	c.JSON(http.StatusOK, MonitorBackward(data))
			// }).Nothing(func() {
			// 	c.JSON(http.StatusNotFound, "")
			// })
		})
}

func (MonitorApi) CreateMonitor(c *gin.Context, request api.PutMonitorRequest) {
	account, ok := c.Value("account").(model.Account)
	if ok {
		// return
	}
	services.CreateMonitor(account, MonitorConfigConvert(request))
}

func (MonitorApi) ListMonitors(c *gin.Context, ordering api.Ordering, index int64, limit int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			// c.JSON(
			// 	http.StatusOK,
			// 	MonitorListBackward(services.ListMonitor(account.ID, index, limit)),
			// )
		})
}

func (MonitorApi) GetMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			// services.GetMonitor(
			// 	account.ID,
			// 	utils.StringToUint(id),
			// ).Just(func(data services.Monitor) {
			// 	c.JSON(
			// 		http.StatusOK,
			// 		// MonitorBackward(data),
			// 	)
			// }).Nothing(func() {
			// 	c.AbortWithStatus(http.StatusNotFound)
			// })
		})
}

func (MonitorApi) DeleteMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			// services.DeleteMonitor(account.ID, utils.StringToUint(id))
			// c.AbortWithStatus(http.StatusNoContent)
		})
}

func (MonitorApi) ListMonitoringRecords(c *gin.Context, id string, index, limit, startAt, endAt int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			// services.GetMonitor(
			// 	account.ID,
			// 	utils.StringToUint(id),
			// ).Just(
			// 	func(data services.Monitor) {
			// 		c.JSON(
			// 			http.StatusOK,
			// 			MonitoringRecordListBackward(services.ListMonitoringRecords(data.Id, index, limit, startAt, endAt)),
			// 		)
			// 	}).Nothing(
			// 	func() {
			// 		c.AbortWithStatus(http.StatusNotFound)
			// 	})
		})
}
