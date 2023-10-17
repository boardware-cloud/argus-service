package argus

import (
	"time"

	"github.com/boardware-cloud/common/constants"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
	"github.com/boardware-cloud/model/notification"
)

type ArgusConfig struct {
	Name                    string        `json:"name"`
	Description             string        `json:"description"`
	Type                    string        `json:"type"`
	Status                  string        `json:"status"`
	MonitorConfig           MonitorConfig `json:"config"`
	NotificationGroupConfig NotificationGroupConfig
}

func (a ArgusConfig) ToEntity() argusModel.Argus {
	argus := argusModel.Argus{
		Type:              constants.MonitorType(a.Type),
		Name:              a.Name,
		Description:       a.Description,
		Status:            constants.MonitorStatus(a.Status),
		NotificationGroup: a.NotificationGroupConfig.ToEntity(),
	}
	argus.SetMonitor(a.MonitorConfig.ToEntity())
	return argus
}

type MonitorConfig interface {
	ToEntity() argusModel.Monitor
}

type HttpMonitorConfig struct {
	Url                 string               `json:"json"`
	Interval            int64                `json:"interval"`
	Timeout             int64                `json:"timeout"`
	Retries             int64                `json:"retries"`
	Headers             []Pair               `json:"headers"`
	AcceptedStatusCodes []string             `json:"acceptedStatusCodes"`
	Method              constants.HttpMehotd `json:"method"`
}

func (config HttpMonitorConfig) ToEntity() argusModel.Monitor {
	headers := common.PairList{}
	for _, header := range config.Headers {
		headers = append(headers, common.Pair{Left: header.Left, Right: header.Right})
	}
	acceptedStatusCodes := common.StringList{}
	for _, code := range config.AcceptedStatusCodes {
		acceptedStatusCodes = append(acceptedStatusCodes, code)
	}
	return &argusModel.HttpMonitor{
		Type:                "HTTP",
		Url:                 config.Url,
		Timeout:             config.Timeout,
		Interval:            time.Duration(config.Interval) * time.Second,
		HttpMethod:          config.Method,
		Headers:             headers,
		AcceptedStatusCodes: acceptedStatusCodes,
	}
}

type NotificationGroupConfig struct {
	Interval      time.Duration `json:"interval"`
	Notifications []NotificationConfig
}

func (config NotificationGroupConfig) ToEntity() notification.NotificationGroup {
	n := notification.NotificationGroup{
		Interval: config.Interval,
	}
	var notificationList []notification.Notification
	for _, c := range config.Notifications {
		notificationList = append(notificationList, c.ToEntity())
	}
	n.SetNotifications(notificationList)
	return n
}

type NotificationConfig interface {
	ToEntity() notification.Notification
}

type EmailNotificationConfig struct {
	Interval  *time.Duration `json:"interval"`
	Receivers EmailReceivers `json:"receivers"`
}

func (e EmailNotificationConfig) ToEntity() notification.Notification {
	n := notification.Notification{
		Type:     "EMAIL",
		Interval: e.Interval,
	}
	n.SetEntity(notification.Email{
		To:  e.Receivers.To,
		Cc:  e.Receivers.Cc,
		Bcc: e.Receivers.Bcc,
	})
	return n
}

type EmailReceivers struct {
	To  []string
	Cc  []string
	Bcc []string
}

type PingMonitorConfig struct {
	Host     string `json:"host"`
	Interval int64  `json:"interval"`
	Retries  int64  `json:"retries"`
}

func (config PingMonitorConfig) ToEntity() argusModel.Monitor {
	return &argusModel.PingMonitor{
		Host:     config.Host,
		Type:     "PING",
		Interval: time.Duration(config.Interval) * time.Second,
	}
}

type Pair struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}
