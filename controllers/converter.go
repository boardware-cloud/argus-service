package controllers

import (
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	model "github.com/boardware-cloud/model/argus"

	api "github.com/boardware-cloud/argus-api"
	services "github.com/boardware-cloud/argus-service/services"
	"github.com/chenyunda218/golambda"
)

func MonitorBackward(monitor services.Monitor) api.Monitor {
	method := api.HttpMethod(*monitor.HttpMethod)
	m := api.Monitor{
		Id:                   utils.UintToString(monitor.Id),
		Name:                 monitor.Name,
		Description:          monitor.Description,
		Type:                 api.MonitorType(monitor.Type),
		Interval:             monitor.Interval,
		Timeout:              monitor.Timeout,
		Url:                  monitor.Url,
		Method:               &method,
		Notifications:        NotificationsBackward(monitor.Notifications),
		NotificationInterval: monitor.NotificationInterval,
		Status:               api.MonitorStatus(monitor.Status),
	}
	return m
}

func MonitorForward(monitor api.Monitor) services.Monitor {
	var retries int64 = 3
	if monitor.Retries <= 10 {
		retries = monitor.Retries
	}
	method := constants.HttpMehotd(*monitor.Method)
	return services.Monitor{
		Id:                   utils.StringToUint(monitor.Id),
		Name:                 monitor.Name,
		Description:          monitor.Description,
		Status:               constants.MonitorStatus(monitor.Status),
		Type:                 constants.MonitorType(monitor.Type),
		Interval:             monitor.Interval,
		Timeout:              monitor.Timeout,
		HttpMethod:           &method,
		Url:                  monitor.Url,
		Notifications:        NotificationsForward(monitor.Notifications),
		NotificationInterval: monitor.NotificationInterval,
		Reties:               retries,
	}
}

func EmailReceiversForwar(receivers api.EmailReceivers) model.EmailReceivers {
	return model.EmailReceivers{
		To:  receivers.To,
		Cc:  receivers.Cc,
		Bcc: receivers.Bcc,
	}
}

func EmailReceiversBackward(receivers model.EmailReceivers) api.EmailReceivers {
	return api.EmailReceivers{
		To:  receivers.To,
		Cc:  receivers.Cc,
		Bcc: receivers.Bcc,
	}
}

func NotificationForward(notification api.Notification) model.Notification {
	receivers := EmailReceiversForwar(*notification.EmailReceivers)
	return model.Notification{
		Type:           constants.EMAIL,
		EmailReceivers: &receivers,
	}
}

func NotificationBackward(notification model.Notification) api.Notification {
	receivers := EmailReceiversBackward(*notification.EmailReceivers)
	return api.Notification{
		Type:           api.NotificationType(notification.Type),
		EmailReceivers: &receivers,
	}
}

func NotificationsForward(notifications []api.Notification) model.Notifications {
	var n model.Notifications
	for _, notification := range notifications {
		n = append(n, NotificationForward(notification))
	}
	return n
}

func NotificationsBackward(notifications model.Notifications) []api.Notification {
	var n []api.Notification
	for _, notification := range notifications {
		n = append(n, NotificationBackward(notification))
	}
	return n
}

func MonitorListBackward(monitorList services.List[services.Monitor]) api.MonitorList {
	list := golambda.Map(monitorList.Data, func(_ int, monitor services.Monitor) api.Monitor {
		return MonitorBackward(monitor)
	})
	return api.MonitorList{
		Data:       list,
		Pagination: PaginationBackward(monitorList),
	}
}

func PaginationBackward[T any](list services.List[T]) api.Pagination {
	return api.Pagination{
		Total: list.Pagination.Total,
		Limit: list.Pagination.Limit,
		Index: list.Pagination.Index,
	}
}

func MonitoringRecordBackward(record services.MonitoringRecord) api.MonitoringRecord {
	id := utils.UintToString(record.Id)
	MonitorId := utils.UintToString(record.MonitorId)
	createdAt := record.CheckedAt.Unix()
	statusCode := record.StatusCode
	result := api.MonitoringResult(record.Result)
	return api.MonitoringRecord{
		Id:           id,
		MonitorId:    MonitorId,
		CheckedAt:    createdAt,
		StatusCode:   statusCode,
		Result:       result,
		ResponseTime: record.ResponseTime,
	}
}

func MonitoringRecordListBackward(recordList services.List[services.MonitoringRecord]) api.MonitoringRecordList {
	data := golambda.Map(recordList.Data, func(_ int, record services.MonitoringRecord) api.MonitoringRecord {
		return MonitoringRecordBackward(record)
	})
	return api.MonitoringRecordList{
		Data:       data,
		Pagination: PaginationBackward(recordList),
	}
}
