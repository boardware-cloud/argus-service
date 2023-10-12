package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	api "github.com/boardware-cloud/argus-api"
	"github.com/boardware-cloud/argus-service/argus"
	"github.com/boardware-cloud/common/config"
	"github.com/boardware-cloud/common/constants"
	argusModel "github.com/boardware-cloud/model/argus"
	"github.com/boardware-cloud/model/common"
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
	}
	return nil
}

func PaginationBackward(p common.Pagination) api.Pagination {
	return api.Pagination{
		Index: p.Index,
		Limit: p.Limit,
		Total: p.Total,
	}
}

func MonitorBackward(a argus.Argus) api.Monitor {
	fmt.Println(a.ID())
	apiModel := api.Monitor{
		Id:          a.ID(),
		Name:        a.Name(),
		Description: a.Description(),
		Type:        api.MonitorType(a.Type()),
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
