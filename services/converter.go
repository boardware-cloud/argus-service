package services

import model "github.com/boardware-cloud/model/argus"

func MonitorBackward(monitor model.Monitor) Monitor {
	return Monitor{
		Id:                   monitor.ID,
		Name:                 monitor.Name,
		Description:          monitor.Description,
		Status:               monitor.Status,
		Url:                  monitor.Url,
		BaseTime:             monitor.BaseTime,
		Timeout:              monitor.Timeout,
		Interval:             monitor.Interval,
		Type:                 monitor.Type,
		HttpMethod:           monitor.HttpMethod,
		UpdatedAt:            monitor.UpdatedAt,
		Notifications:        monitor.Notifications,
		NotificationInterval: monitor.NotificationInterval,
		Reties:               monitor.Retries,
	}
}

func UptimeNodeForward(node UptimeNode) model.UptimeNode {
	m := model.UptimeNode{
		Heartbeat:         node.Heartbeat,
		HeartbeatInterval: node.HeartbeatInterval,
	}
	m.ID = node.ID
	return m
}

func MonitoringResultBackward(m model.MonitoringRecord) MonitoringRecord {
	return MonitoringRecord{
		Id:           m.ID,
		MonitorId:    m.MonitorId,
		Result:       m.Result,
		CheckedAt:    m.CheckedAt,
		StatusCode:   m.StatusCode,
		ResponseTime: m.ResponseTime,
	}
}
