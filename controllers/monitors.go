package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/common/utils"
	model "github.com/boardware-cloud/core-api"
	"github.com/gin-gonic/gin"
)

const DEFAULT_TIMEOUT = 10
const DEFAULT_INTERVAL = 5 * 60

type MonitorApi struct{}

func (MonitorApi) UpdateMonitor(c *gin.Context, monitorId string, putMonitorRequest api.PutMonitorRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			services.UpdateMonitor(
				utils.StringToUint(account.Id),
				utils.StringToUint(monitorId),
				PutMonitorForward(putMonitorRequest),
			).Just(func(data services.Monitor) {
				c.JSON(http.StatusOK, MonitorBackward(data))
			}).Nothing(func() {
				c.JSON(http.StatusNotFound, "")
			})
		})
}

func (MonitorApi) CreateMonitor(c *gin.Context, createMonitorRequest api.PutMonitorRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			c.JSON(
				http.StatusCreated,
				MonitorBackward(services.CreateMonitor(
					utils.StringToUint(account.Id), PutMonitorForward(createMonitorRequest),
				)))
		})
}

func (MonitorApi) ListMonitors(c *gin.Context, ordering api.Ordering, index int64, limit int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			c.JSON(
				http.StatusOK,
				MonitorListBackward(services.ListMonitor(utils.StringToUint(account.Id), index, limit)),
			)
		})
}

func (MonitorApi) GetMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			services.GetMonitor(
				utils.StringToUint(account.Id),
				utils.StringToUint(id),
			).Just(func(data services.Monitor) {
				c.JSON(
					http.StatusOK,
					MonitorBackward(data),
				)
			}).Nothing(func() {
				c.AbortWithStatus(http.StatusNotFound)
			})
		})
}

func (MonitorApi) DeleteMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			services.DeleteMonitor(utils.StringToUint(account.Id), utils.StringToUint(id))
			c.AbortWithStatus(http.StatusNoContent)
		})
}

func (MonitorApi) ListMonitoringRecords(c *gin.Context, id string, index, limit, startAt, endAt int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			services.GetMonitor(
				utils.StringToUint(account.Id),
				utils.StringToUint(id),
			).Just(
				func(data services.Monitor) {
					c.JSON(
						http.StatusOK,
						MonitoringRecordListBackward(services.ListMonitoringRecords(data.Id, index, limit, startAt, endAt)),
					)
				}).Nothing(
				func() {
					c.AbortWithStatus(http.StatusNotFound)
				})
		})
}
