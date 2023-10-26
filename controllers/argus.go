package controllers

import (
	"net/http"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/services"
	"github.com/boardware-cloud/common/code"
	"github.com/boardware-cloud/common/utils"
	coreServices "github.com/boardware-cloud/core/services"
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
			a, err := argusService.UpdateMonitor(utils.StringToUint(monitorId), MonitorConfigConvert(request))
			if err != nil {
				code.GinHandler(c, err)
				return
			}
			c.JSON(http.StatusOK, MonitorBackward(a))
		})
}

func (MonitorApi) CreateMonitor(ctx *gin.Context, request api.PutMonitorRequest) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	a, err := argusService.CreateMonitor(account.Entity, MonitorConfigConvert(request))
	if err != nil {
		code.GinHandler(ctx, err)
	}
	ctx.JSON(http.StatusOK, MonitorBackward(a))
}

func (MonitorApi) ListMonitors(ctx *gin.Context, ordering api.Ordering, index int64, limit int64) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	list, pagination := argusService.ListMonitors(account.ID(), index, limit)
	var monitorList []api.Monitor
	for _, item := range list {
		monitorList = append(monitorList, Convert(item).(api.Monitor))
	}
	ctx.JSON(http.StatusOK, api.MonitorList{
		Data:       monitorList,
		Pagination: Convert(pagination).(api.Pagination),
	})
}

func (MonitorApi) GetMonitor(ctx *gin.Context, id string) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	a, err := argusService.GetMonitor(utils.StringToUint(id))
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}
	if !account.Entity.Own(&a) {
		code.GinHandler(ctx, code.ErrPermissionDenied)
		return
	}
	ctx.JSON(http.StatusOK, MonitorBackward(a))
}

func (MonitorApi) DeleteMonitor(ctx *gin.Context, id string) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	monitor, err := argusService.GetMonitor((utils.StringToUint(id)))
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}
	if !account.Entity.Own(&monitor) {
		code.GinHandler(ctx, code.ErrPermissionDenied)
		return
	}
	services.DeleteMonitor(monitor)
	ctx.JSON(http.StatusNoContent, "")
}

func (MonitorApi) ListMonitoringRecords(ctx *gin.Context, id string, index, limit, startAt, endAt int64) {
	account := coreServices.GetAccount(ctx)
	if account == nil {
		return
	}
	monitor, err := argusService.GetMonitor((utils.StringToUint(id)))
	if err != nil {
		code.GinHandler(ctx, err)
		return
	}
	if !account.Entity.Own(&monitor) {
		code.GinHandler(ctx, code.ErrPermissionDenied)
		return
	}
	list, pagination := services.ListRecords(monitor.Entity().ID, index, limit)
	var recordList []api.MonitoringRecord
	for _, record := range list {
		recordList = append(recordList, Convert(record).(api.MonitoringRecord))
	}
	ctx.JSON(http.StatusOK, api.MonitoringRecordList{
		Data:       recordList,
		Pagination: Convert(pagination).(api.Pagination),
	})
}
