package argus

import (
	"github.com/boardware-cloud/common/constants"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/core"
)

type ArgusConfig struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Type          string        `json:"type"`
	Status        string        `json:"status"`
	MonitorConfig MonitorConfig `json:"config"`
}

func (a *ArgusConfig) FromEntity(argus argusModel.Argus) ArgusConfig {
	a.Name = argus.Name
	a.Description = argus.Description
	a.Type = string(argus.Type)
	a.Status = string(argus.Status)
	return *a
}

func (a ArgusConfig) ToEntity(account core.Account) argusModel.Argus {
	argus := argusModel.Argus{
		Type:        constants.MonitorType(a.Type),
		Name:        a.Name,
		Description: a.Description,
		Status:      constants.MonitorStatus(a.Status),
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
	Headers             Pair                 `json:"headers"`
	AcceptedStatusCodes []string             `json:"acceptedStatusCodes"`
	Method              constants.HttpMehotd `json:"method"`
}

func (config HttpMonitorConfig) ToEntity() argusModel.Monitor {
	return &argusModel.HttpMonitor{
		Type:       "HTTP",
		Url:        config.Url,
		Timeout:    config.Timeout,
		Interval:   config.Interval,
		HttpMethod: config.Method,
	}
}

type PingMonitorConfig struct {
	Url      string `json:"json"`
	Interval int64  `json:"interval"`
	Timeout  int64  `json:"timeout"`
	Retries  int64  `json:"retries"`
}

func (config PingMonitorConfig) ToEntity() argusModel.Monitor {
	return &argusModel.PingMonitor{
		Type:     "PING",
		Url:      config.Url,
		Timeout:  config.Timeout,
		Interval: config.Interval,
	}
}

type Pair struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}
