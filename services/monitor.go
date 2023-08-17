package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boardware-cloud/common/constants"
	model "github.com/boardware-cloud/model/argus"
	f "github.com/chenyunda218/golambda"
)

type Monitor struct {
	Id                   uint
	Name                 string
	Description          string
	Status               constants.MonitorStatus
	Type                 constants.MonitorType
	Interval             int64
	Timeout              int64
	HttpMethod           *constants.HttpMehotd
	BaseTime             int64
	Url                  string
	Heartbeat            int64
	UpdatedAt            time.Time
	Notifications        model.Notifications
	NotificationInterval int64
	Reties               int64
}

type MonitoringRecord struct {
	Id           uint
	MonitorId    uint
	CheckedAt    time.Time
	Result       constants.MonitoringResult
	ResponseTime *int64
	StatusCode   string
}

func (m *Monitor) httpMonitor() model.MonitoringRecord {
	client := &http.Client{Timeout: time.Duration(m.Timeout) * time.Second}
	req, _ := http.NewRequest(string(*m.HttpMethod), m.Url, nil)
	start := time.Now().UnixMilli()
	resp, err := client.Do(req)
	checkedAt := time.Now()
	responseTime := time.Now().UnixMilli() - start
	record := model.MonitoringRecord{
		MonitorId:    m.Id,
		CheckedAt:    checkedAt,
		Url:          m.Url,
		Type:         m.Type,
		HttpMethod:   m.HttpMethod,
		ResponseTime: &responseTime,
	}
	if err != nil {
		if resp == nil {
			record.Result = constants.TIMEOUT
		} else {
			record.Result = constants.DOWN
		}
	} else {
		record.StatusCode = fmt.Sprint(resp.StatusCode)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			record.Result = constants.OK
		} else {
			record.Result = constants.DOWN
		}
	}
	return record
}

func (m *Monitor) Notification(record model.MonitoringRecord) {
	LastAlert := &model.UptimeMonitorAlert{}
	ctx := DB.Where("monitor_id = ?", m.Id).Limit(1).Order("created_at DESC").Find(LastAlert)
	if ctx.RowsAffected == 0 || LastAlert.CreatedAt.Unix()+m.NotificationInterval <= time.Now().Unix() {
		alert := &model.UptimeMonitorAlert{}
		alert.MonitorId = m.Id
		alert.Notifications = m.Notifications
		alert.Subject = "BoardWare Uptime Monitor alert"
		alert.Message = fmt.Sprintf(`
		<html>
			<body>
				<div>Url: %s</div>
				<div>Check time: %s</div>
				<div>Status: %s</div>
			</body>
		</html>`, m.Url, record.CheckedAt, record.Result)
		DB.Save(&alert)
		for _, notifiction := range m.Notifications {
			switch notifiction.Type {
			case constants.EMAIL:
				emailSender.SendHtml(
					"dan.chen@boardware.com",
					alert.Subject,
					alert.Message,
					notifiction.EmailReceivers.To,
					notifiction.EmailReceivers.Cc,
					notifiction.EmailReceivers.Bcc,
				)
			}
		}
	}
}

func (m *Monitor) Check() bool {
	var retries int64 = 0
	for {
		now := time.Now()
		m.Heartbeat = now.Unix()
		currentMonitor := GetMonitorById(m.Id).Value()
		if currentMonitor == nil || m.UpdatedAt != currentMonitor.UpdatedAt || currentMonitor.Status != constants.ACTIVED {
			return false
		}
		m = currentMonitor
		var record model.MonitoringRecord
		switch m.Type {
		case constants.HTTP:
			record = m.httpMonitor()
		}
		DB.Save(&record)
		if record.Result == constants.OK {
			break
		}
		if retries >= m.Reties {
			m.Notification(record)
			break
		}
		retries++
	}
	return true
}

func CreateMonitor(
	accountId uint,
	monitor model.Monitor,
) Monitor {
	monitor.AccountId = accountId
	monitor.Retries = 3
	DB.Save(&monitor)
	m := MonitorBackward(monitor)
	go Spawn(m)
	return m
}

func OrphanMonitor() []Monitor {
	var monitors []model.Monitor
	DB.Model(&model.Monitor{}).Where(
		"uptime_node_id IS NULL AND status = ? AND deleted_at IS NULL AND status = 'ACTIVED'",
		constants.ACTIVED,
	).Find(&monitors)
	return f.Map(monitors, func(_ int, monitor model.Monitor) Monitor {
		return MonitorBackward(monitor)
	})
}

func ListMonitor(accountId uint, index, limit int64) List[Monitor] {
	var monitors []model.Monitor
	var total int64
	ctx := DB.Model(&model.Monitor{}).Where("account_id = ?", accountId).Count(&total)
	if total == 0 {
		return List[Monitor]{
			Data: []Monitor{},
			Pagination: Pagination{
				Limit: 1,
				Index: 0,
				Total: 0,
			},
		}
	}
	if total <= index*limit {
		index = total/limit - 1
	}
	ctx.Limit(int(limit)).Offset(int(index * limit)).Find(&monitors)
	return List[Monitor]{
		Data: f.Map(monitors, func(_ int, monitor model.Monitor) Monitor {
			return MonitorBackward(monitor)
		}),
		Pagination: Pagination{
			Limit: limit,
			Index: index,
			Total: total,
		},
	}
}

func GetMonitor(accountId uint, monitorId uint) f.MayBe[Monitor] {
	monitor := model.Monitor{
		AccountId: accountId,
	}
	monitor.ID = monitorId
	ctx := DB.Where("account_id = ?", accountId).Find(&monitor)
	if ctx.RowsAffected == 0 {
		return f.MayBe[Monitor]{}
	}
	m := MonitorBackward(monitor)
	return f.MayBe[Monitor]{
		Data: &m,
	}
}

func GetMonitorById(monitorId uint) f.MayBe[Monitor] {
	monitor := model.Monitor{}
	ctx := DB.Find(&monitor, monitorId)
	if ctx.RowsAffected == 0 {
		return f.MayBe[Monitor]{}
	}
	m := MonitorBackward(monitor)
	return f.MayBe[Monitor]{Data: &m}
}

func ListMonitoringRecords(monitorId uint, index, limit, startAt, endAt int64) List[MonitoringRecord] {
	var records []model.MonitoringRecord
	var total int64
	ctx := DB.Model(&model.MonitoringRecord{}).Where("monitor_id = ?", monitorId)
	if startAt > endAt && endAt != 0 {
		temp := endAt
		endAt = startAt
		startAt = temp
	}
	if startAt != 0 {
		ctx = ctx.Where("checked_at >= ?", time.Unix(startAt, 0))
	}
	if endAt != 0 {
		fmt.Println(endAt)
		ctx = ctx.Where("checked_at < ?", time.Unix(endAt+1, 0))
	}
	ctx.Count(&total)
	if total <= index*limit {
		index = total/limit - 1
		if index < 0 {
			index = 0
		}
	}
	ctx.Order("checked_at DESC").Limit(int(limit)).Offset(int(index * limit)).Find(&records)
	return List[MonitoringRecord]{
		Data: f.Map(records, func(_ int, m model.MonitoringRecord) MonitoringRecord {
			return MonitoringResultBackward(m)
		}),
		Pagination: Pagination{
			Limit: limit,
			Index: index,
			Total: total,
		},
	}
}

func DeleteMonitor(accountId, monitorId uint) {
	monitor := &model.Monitor{
		AccountId: accountId,
	}
	monitor.ID = monitorId
	DB.Delete(&monitor)
}

func UpdateMonitor(
	accountId,
	monitorId uint,
	name string,
	description string,
	Type constants.MonitorType,
	interval int64,
	timeout int64,
	retries int64,
	httpMethod *constants.HttpMehotd,
	url string,
	notifications model.Notifications,
	notificationInterval int64,
	status constants.MonitorStatus,
) f.MayBe[Monitor] {
	var monitor model.Monitor
	monitor.ID = monitorId
	monitor.AccountId = accountId
	ctx := DB.Where("account_id = ?", accountId).Find(&monitor, monitorId)
	if ctx.RowsAffected == 0 {
		return f.MayBe[Monitor]{}
	}
	monitor.Name = name
	monitor.Description = description
	monitor.Type = Type
	monitor.Interval = interval
	monitor.Timeout = timeout
	monitor.Retries = retries
	monitor.HttpMethod = httpMethod
	monitor.Url = url
	monitor.Notifications = notifications
	monitor.NotificationInterval = notificationInterval
	monitor.UptimeNodeId = nil
	monitor.Status = status
	monitor.Retries = 3
	DB.Save(&monitor)
	m := MonitorBackward(monitor)
	go Spawn(m)
	return f.MayBe[Monitor]{
		Data: &m,
	}
}
