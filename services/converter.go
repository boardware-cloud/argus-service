package services

import (
	model "github.com/boardware-cloud/model/argus"
	common "github.com/boardware-cloud/model/common"
	f "github.com/chenyunda218/golambda"
)

func PairListBackward(pairList *common.PairList) *[]Pair {
	if pairList == nil {
		return nil
	}
	return f.Reference(f.Map(*pairList, func(_ int, pair common.Pair) Pair {
		return Pair{
			Left:  pair.Left,
			Right: pair.Right,
		}
	}))
}

func MonitorBackward(monitor model.Monitor) Monitor {
	var acceptedStatusCodes *[]string
	if monitor.AcceptedStatusCodes != nil {
		acceptedStatusCodes = (*[]string)(monitor.AcceptedStatusCodes)
	}
	return Monitor{
		Id:                   monitor.ID,
		Name:                 monitor.Name,
		Description:          monitor.Description,
		Status:               monitor.Status,
		Url:                  monitor.Url,
		Timeout:              monitor.Timeout,
		Interval:             monitor.Interval,
		Type:                 monitor.Type,
		HttpMethod:           monitor.HttpMethod,
		UpdatedAt:            monitor.UpdatedAt,
		Notifications:        monitor.Notifications,
		NotificationInterval: monitor.NotificationInterval,
		Reties:               monitor.Retries,
		Body:                 monitor.Body,
		Headers:              PairListBackward(monitor.Headers),
		AcceptedStatusCodes:  acceptedStatusCodes,
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
		Body:         m.Body,
	}
}
