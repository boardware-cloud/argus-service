package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/utils"
	"github.com/boardware-cloud/middleware"
	model "github.com/boardware-cloud/model/core"
	"github.com/gin-gonic/gin"
)

const DEFAULT_TIMEOUT = 10
const DEFAULT_INTERVAL = 5 * 60

type MonitorApi struct{}

func (MonitorApi) UpdateMonitor(c *gin.Context, monitorId string, request api.PutMonitorRequest) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			a, err := services.UpdateMonitor(utils.StringToUint(monitorId), MonitorConfigConvert(request))
			if err != nil {
				code.GinHandler(c, err)
				return
			}
			c.JSON(http.StatusOK, MonitorBackward(a))
		})
}

func (MonitorApi) CreateMonitor(ctx *gin.Context, request api.PutMonitorRequest) {
	middleware.GetAccount(ctx, func(c *gin.Context, account model.Account) {
		a, err := services.CreateMonitor(account, MonitorConfigConvert(request))
		if err != nil {
			code.GinHandler(c, err)
		}
		ctx.JSON(http.StatusOK, MonitorBackward(a))
	})
}

func (MonitorApi) ListMonitors(c *gin.Context, ordering api.Ordering, index int64, limit int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			list, pagination := services.ListMonitors(account.ID(), index, limit)
			var monitorList []api.Monitor
			for _, item := range list {
				monitorList = append(monitorList, Convert(item).(api.Monitor))
			}
			c.JSON(http.StatusOK, api.MonitorList{
				Data:       monitorList,
				Pagination: Convert(pagination).(api.Pagination),
			})
		})
}

func (MonitorApi) GetMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			a, err := services.GetMonitor(utils.StringToUint(id))
			if err != nil {
				code.GinHandler(c, err)
				return
			}
			if !account.Own(&a) {
				code.GinHandler(c, code.ErrPermissionDenied)
				return
			}
			c.JSON(http.StatusOK, MonitorBackward(a))
		})
}

func (MonitorApi) DeleteMonitor(c *gin.Context, id string) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			monitor, err := services.GetMonitor((utils.StringToUint(id)))
			if err != nil {
				code.GinHandler(c, err)
				return
			}
			if !account.Own(&monitor) {
				code.GinHandler(c, code.ErrPermissionDenied)
				return
			}
			services.DeleteMonitor(monitor)
			c.JSON(http.StatusNoContent, "")
		})
}

func (MonitorApi) ListMonitoringRecords(c *gin.Context, id string, index, limit, startAt, endAt int64) {
	middleware.GetAccount(c,
		func(c *gin.Context, account model.Account) {
			monitor, err := services.GetMonitor((utils.StringToUint(id)))
			if err != nil {
				code.GinHandler(c, err)
				return
			}
			if !account.Own(&monitor) {
				code.GinHandler(c, code.ErrPermissionDenied)
				return
			}
			list, pagination := services.ListRecords(monitor.Entity().ID, index, limit)
			var recordList []api.MonitoringRecord
			for _, item := range list {
				recordList = append(recordList, Convert(item).(api.MonitoringRecord))
			}
			c.JSON(http.StatusOK, api.MonitoringRecordList{
				Data:       recordList,
				Pagination: Convert(pagination).(api.Pagination),
			})
		})
}
