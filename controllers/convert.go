package controllers

import (
	"encoding/json"
	"time"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/notification"
	"github.com/chenyunda218/golambda"
)

func Convert(from any) any {
	switch f := from.(type) {
	case api.PutMonitorRequest:
		return MonitorConfigConvert(f)
	case api.PingMonitor:
		return PingMonitorConfigConvert(f)
	case api.HttpMonitor:
		return HttpMonitorConfigConvert(f)
	case argus.Argus:
		return MonitorBackward(f)
	case common.Pagination:
		return PaginationBackward(f)
	case argus.Record:
		return RecordBackward(f)
	}
	return nil
}

func RecordBackward(a argus.Record) api.MonitoringRecord {
	return api.MonitoringRecord{
		ResponseTime: int64(a.ResponesTime) / int64(time.Second),
		Result:       api.MonitoringResult(a.Result),
		CheckedAt:    a.CheckedAt.Unix(),
	}
}

func PaginationBackward(p common.Pagination) api.Pagination {
	return api.Pagination{
		Index: p.Index,
		Limit: p.Limit,
		Total: p.Total,
	}
}

func MonitorBackward(a argus.Argus) api.Monitor {
	notificationGroup := NotificationGroupBackward(a.Entity().NotificationGroup)
	apiModel := api.Monitor{
		Status:            api.MonitorStatus(a.Entity().Status),
		Id:                utils.UintToString(a.ID()),
		Name:              a.Name(),
		Description:       a.Description(),
		Type:              api.MonitorType(a.Type()),
		NotificationGroup: &notificationGroup,
	}
	switch argusMonitor := a.Monitor().(type) {
	case *argus.HttpMonitor:
		m := argusMonitor.Entity().(*argusModel.HttpMonitor)
		j, _ := json.Marshal(m)
		apiHttp := api.HttpMonitor{}
		json.Unmarshal(j, &apiHttp)
		interval := *apiHttp.Interval / int64(time.Second)
		apiHttp.Interval = &interval
		apiModel.HttpMonitor = &apiHttp
	case *argus.PingMonitor:
		m := argusMonitor.Entity().(*argusModel.PingMonitor)
		j, _ := json.Marshal(m)
		apiPing := api.PingMonitor{}
		json.Unmarshal(j, &apiPing)
		apiModel.PingMonitor = &apiPing
	}
	return apiModel
}

func MonitorConfigConvert(raw api.PutMonitorRequest) argus.ArgusConfig {
	var monitorConfig argus.MonitorConfig
	t := config.Convention(raw.Type, api.HTTP)
	switch config.Convention(raw.Type, api.HTTP) {
	case api.HTTP:
		monitorConfig = HttpMonitorConfigConvert(config.Convention(raw.HttpMonitor, api.HttpMonitor{}))
	case api.PING:
		monitorConfig = PingMonitorConfigConvert(config.Convention(raw.PingMonitor, api.PingMonitor{}))
	}
	return argus.ArgusConfig{
		Name:                    config.Convention(raw.Name, ""),
		Description:             config.Convention(raw.Description, ""),
		Status:                  string(config.Convention(raw.Status, api.DISACTIVED)),
		Type:                    string(t),
		MonitorConfig:           monitorConfig,
		NotificationGroupConfig: NotificationGroupConvert(config.Convention(raw.NotificationGroup, api.NotificationGroup{})),
	}
}

func PingMonitorConfigConvert(raw api.PingMonitor) argus.PingMonitorConfig {
	return argus.PingMonitorConfig{
		Url:      config.Convention(raw.Url, ""),
		Interval: config.Convention(raw.Interval, 60),
		Timeout:  config.Convention(raw.Timeout, 10),
		Retries:  config.Convention(raw.Retries, 3),
	}
}

func HttpMonitorConfigConvert(raw api.HttpMonitor) argus.HttpMonitorConfig {
	return argus.HttpMonitorConfig{
		Url:      config.Convention(raw.Url, ""),
		Interval: config.Convention(raw.Interval, 60),
		Timeout:  config.Convention(raw.Timeout, 10),
		Retries:  config.Convention(raw.Retries, 3),
		Method:   constants.HttpMehotd(config.Convention(raw.Method, api.GET)),
		Headers: golambda.Map(config.Convention(raw.Headers, []api.Pair{}), func(_ int, header api.Pair) argus.Pair {
			return argus.Pair{
				Left:  header.Left,
				Right: header.Right,
			}
		}),
		AcceptedStatusCodes: config.Convention(raw.AcceptedStatusCodes, []string{}),
	}
}

func NotificationGroupBackward(raw notification.NotificationGroup) api.NotificationGroup {
	interval := int64(raw.Interval / time.Second)
	var notifications []api.Notification
	for _, n := range raw.Notifications() {
		temp := api.Notification{
			Interval: &interval,
			Type:     api.NotificationType(n.Type),
		}
		switch n.Type {
		case "EMAIL":
			entity := n.Entity().(notification.Email)
			temp.Email = &api.EmailNotification{
				Receivers: &api.EmailReceivers{
					To:  entity.To,
					Cc:  entity.Cc,
					Bcc: entity.Bcc,
				},
			}
		}
		notifications = append(notifications, temp)
	}
	o := api.NotificationGroup{
		Interval:      &interval,
		Notifications: &notifications,
	}
	return o
}

func NotificationGroupConvert(raw api.NotificationGroup) argus.NotificationGroupConfig {
	return argus.NotificationGroupConfig{
		Interval: time.Second * time.Duration(config.Convention(raw.Interval, int64(600))),
		Notifications: golambda.Map(
			config.Convention(
				raw.Notifications,
				[]api.Notification{},
			),
			func(_ int, a api.Notification) argus.NotificationConfig {
				return NotificationConvert(a)
			}),
	}
}

func NotificationConvert(raw api.Notification) argus.NotificationConfig {
	return EmailNotificationConvert(config.Convention(raw.Email, api.EmailNotification{}))
}

func EmailNotificationConvert(raw api.EmailNotification) argus.EmailNotificationConfig {
	return argus.EmailNotificationConfig{
		Receivers: argus.EmailReceivers{
			To:  raw.Receivers.To,
			Cc:  raw.Receivers.Cc,
			Bcc: raw.Receivers.Bcc,
		},
	}
}
