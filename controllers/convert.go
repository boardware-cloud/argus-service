package controllers

import (
	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/constants"
)

func MonitorConfigConvert(raw api.PutMonitorRequest) argus.ArgusConfig {
	var monitorConfig argus.MonitorConfig
	t := config.Convention(raw.Type, api.HTTP)
	switch t {
	case api.HTTP:
		monitorConfig = HttpMonitorConfigConvert(config.Convention(raw.HttpMonitor, api.HttpMonitor{}))
	case api.PING:
		monitorConfig = PingMonitorConfigConvert(config.Convention(raw.PingMonitor, api.PingMonitor{}))
	}
	return argus.ArgusConfig{
		Name:          config.Convention(raw.Name, ""),
		Description:   config.Convention(raw.Description, ""),
		Status:        string(config.Convention(raw.Status, api.DISACTIVED)),
		Type:          string(t),
		MonitorConfig: monitorConfig,
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
	}
}
