package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	model "github.com/boardware-cloud/core-api"
	uptime "github.com/boardware-cloud/model/argus"
	f "github.com/chenyunda218/golambda"
	"github.com/gin-gonic/gin"
)

const DEFAULT_TIMEOUT = 10
const DEFAULT_INTERVAL = 5 * 60

type MonitorApi struct{}

func (MonitorApi) UpdateMonitor(c *gin.Context, monitorId string, updateMonitorRequest api.PutMonitorRequest) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		method := constants.HttpMehotd(*updateMonitorRequest.Method)
		services.UpdateMonitor(
			utils.StringToUint(account.Id),
			utils.StringToUint(monitorId),
			updateMonitorRequest.Name,
			updateMonitorRequest.Description,
			constants.HTTP,
			updateMonitorRequest.Interval,
			updateMonitorRequest.Timeout, 0,
			&method,
			updateMonitorRequest.Url,
			NotificationsForward(updateMonitorRequest.Notifications),
			updateMonitorRequest.NotificationInterval,
			constants.MonitorStatus(updateMonitorRequest.Status),
		).Just(func(data services.Monitor) {
			c.JSON(http.StatusOK, MonitorBackward(data))
		}).Nothing(func() {
			c.JSON(http.StatusNotFound, gin.H{})
		})
	})
}

func (MonitorApi) CreateMonitor(c *gin.Context, createMonitorRequest api.PutMonitorRequest) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		var httpMehtod *constants.HttpMehotd
		f.NewMayBe(createMonitorRequest.Method).Just(func(method api.HttpMethod) {
			httpMehtod = f.Reference(constants.HttpMehotd(method))
		})
		c.JSON(
			http.StatusCreated,
			MonitorBackward(services.CreateMonitor(
				utils.StringToUint(account.Id), uptime.Monitor{
					AccountId:            utils.StringToUint(account.Id),
					Name:                 createMonitorRequest.Name,
					Description:          createMonitorRequest.Description,
					Url:                  createMonitorRequest.Url,
					Status:               constants.MonitorStatus(createMonitorRequest.Status),
					Interval:             createMonitorRequest.Interval,
					Timeout:              createMonitorRequest.Timeout,
					Notifications:        NotificationsForward(createMonitorRequest.Notifications),
					Retries:              0,
					Type:                 constants.MonitorType(createMonitorRequest.Type),
					HttpMethod:           httpMehtod,
					NotificationInterval: createMonitorRequest.NotificationInterval,
				},
			)))
	})
}

func (MonitorApi) ListMonitors(c *gin.Context, ordering api.Ordering, index int64, limit int64) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		c.JSON(
			http.StatusOK,
			MonitorListBackward(services.ListMonitor(utils.StringToUint(account.Id), index, limit)),
		)
	})
}

func (MonitorApi) GetMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
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
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		services.DeleteMonitor(utils.StringToUint(account.Id), utils.StringToUint(id))
		c.AbortWithStatus(http.StatusNoContent)
	})
}

func (MonitorApi) ListMonitoringRecords(c *gin.Context, id string, index, limit, startAt, endAt int64) {
	middleware.GetAccount(c, func(c *gin.Context, account model.Account) {
		services.GetMonitor(
			utils.StringToUint(account.Id),
			utils.StringToUint(id),
		).Just(func(data services.Monitor) {
			c.JSON(
				http.StatusOK,
				MonitoringRecordListBackward(services.ListMonitoringRecords(data.Id, index, limit, startAt, endAt)),
			)
		}).Nothing(func() {
			c.AbortWithStatus(http.StatusNotFound)
		})
	})
}
