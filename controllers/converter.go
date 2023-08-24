package controllers

import (
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	model "github.com/boardware-cloud/model/argus"
	common "github.com/boardware-cloud/model/common"

	api "github.com/boardware-cloud/argus-api"
	services "github.com/boardware-cloud/argus-service/services"
	f "github.com/chenyunda218/golambda"
)

func MonitorBackward(monitor services.Monitor) api.Monitor {
	method := api.HttpMethod(*monitor.HttpMethod)
	var headers *[]api.Pair
	if monitor.Headers != nil {
		var hs []api.Pair
		for _, header := range *monitor.Headers {
			hs = append(hs, api.Pair{
				Left:  header.Left,
				Right: header.Right,
			})
		}
		headers = &hs
	}
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
		Retries:              monitor.Reties,
		Headers:              headers,
		AcceptedStatusCodes:  monitor.AcceptedStatusCodes,
	}
	return m
}

func PairListForward(pairList *[]api.Pair) *common.PairList {
	if pairList == nil {
		return nil
	}
	var pairs common.PairList = make(common.PairList, 0)
	for _, pair := range *pairList {
		pairs = append(pairs, common.Pair{
			Left:  pair.Left,
			Right: pair.Right,
		})
	}
	return &pairs
}

func StringListForward(ss *[]string) *common.StringList {
	if ss == nil {
		return nil
	}
	var list common.StringList = make(common.StringList, 0)
	for _, s := range *ss {
		list = append(list, s)
	}
	return &list
}

func PairListBackward(pairList *common.PairList) *[]api.Pair {
	if pairList == nil {
		return nil
	}
	var pairs []api.Pair = make([]api.Pair, 0)
	for _, pair := range *pairList {
		pairs = append(pairs, api.Pair{
			Left:  pair.Left,
			Right: pair.Right,
		})
	}
	return &pairs
}

func PutMonitorForward(putMonitorRequest api.PutMonitorRequest) model.Monitor {
	var httpMehtod *constants.HttpMehotd
	f.NewMayBe(putMonitorRequest.Method).Just(func(method api.HttpMethod) {
		httpMehtod = f.Reference(constants.HttpMehotd(method))
	})
	var retries int64 = 3
	if putMonitorRequest.Retries <= 10 {
		retries = putMonitorRequest.Retries
	}
	return model.Monitor{
		Name:                 putMonitorRequest.Name,
		Description:          putMonitorRequest.Description,
		Url:                  putMonitorRequest.Url,
		Status:               constants.MonitorStatus(putMonitorRequest.Status),
		Interval:             putMonitorRequest.Interval,
		Timeout:              putMonitorRequest.Timeout,
		Notifications:        NotificationsForward(putMonitorRequest.Notifications),
		Retries:              retries,
		Type:                 constants.MonitorType(putMonitorRequest.Type),
		HttpMethod:           httpMehtod,
		NotificationInterval: putMonitorRequest.NotificationInterval,
		Body:                 nil,
		Headers:              PairListForward(putMonitorRequest.Headers),
		AcceptedStatusCodes:  StringListForward(putMonitorRequest.AcceptedStatusCodes),
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
	list := f.Map(monitorList.Data, func(_ int, monitor services.Monitor) api.Monitor {
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
	data := f.Map(recordList.Data, func(_ int, record services.MonitoringRecord) api.MonitoringRecord {
		return MonitoringRecordBackward(record)
	})
	return api.MonitoringRecordList{
		Data:       data,
		Pagination: PaginationBackward(recordList),
	}
}
