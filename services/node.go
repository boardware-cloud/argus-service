package services

import (
	"time"

	model "github.com/boardware-cloud/model/argus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	HEARTBEAT_TOLERANCE     = 5
	CHECK_MEMBER_INTERVAL   = 10
	HEARTBEAT_INTERVAL      = 10
	CHECK_MONITORS_INTERVAL = 10
)

type UptimeNode struct {
	ID                uint
	Heartbeat         int64
	HeartbeatInterval int64
	Entity            model.UptimeNode
}

func NewUptimeNode() UptimeNode {
	return UptimeNode{}
}

func Spawn(monitor Monitor) {
	tx := DB.Model(&model.Monitor{}).Where(
		"id = ? AND uptime_node_id IS NULL AND status = 'ACTIVED'",
		monitor.Id).Update("uptime_node_id", node.ID)
	m := GetMonitorById(monitor.Id)
	if m.Data == nil {
		return
	}
	monitor.UpdatedAt = m.Data.UpdatedAt
	records := ListMonitoringRecords(monitor.Id, 0, 1, 0, 0)
	if len(records.Data) > 0 {
		record := records.Data[0]
		time.Sleep(time.Duration(record.CheckedAt.Unix()+monitor.Interval-time.Now().Unix()) * time.Second)
	}
	if tx.RowsAffected != 0 {
		for monitor.Check() {
			time.Sleep(time.Duration(monitor.Interval) * time.Second)
		}
	}
}

func RecoverNode(id uint) {
	DB.Transaction(func(tx *gorm.DB) error {
		node := model.UptimeNode{}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&node, id); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&model.Monitor{}).Where(
			"uptime_node_id = ?",
			node.ID).Update("uptime_node_id", nil); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Delete(&model.UptimeNode{}, id); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		return nil
	})
}

func CheckMembers() {
	for _, node := range ListUptimeNodes() {
		if time.Now().Unix() > node.Heartbeat+node.HeartbeatInterval+HEARTBEAT_TOLERANCE {
			go RecoverNode(node.ID)
		}
	}
}

func CheckMontiors() {
	for _, monitor := range OrphanMonitor() {
		go Spawn(monitor)
	}
}

func (node *UptimeNode) Register() {
	node.HeartbeatInterval = HEARTBEAT_INTERVAL
	node.Heartbeat = time.Now().Unix()
	m := UptimeNodeForward(*node)
	DB.Save(&m)
	DB.Find(&m)
	node.ID = m.ID
	node.Entity = m
}

func (node *UptimeNode) Beat() {
	node.Heartbeat = time.Now().Unix()
	node.Entity.Heartbeat = node.Heartbeat
	DB.Save(&node.Entity)
}

func KeepAlive() {
	for {
		node.Beat()
		time.Sleep(HEARTBEAT_INTERVAL * time.Second)
	}
}

func KeepCheckMontiors() {
	for {
		CheckMontiors()
		time.Sleep(CHECK_MONITORS_INTERVAL * time.Second)
	}
}

func KeepCheckNodes() {
	for {
		CheckMembers()
		time.Sleep(CHECK_MEMBER_INTERVAL * time.Second)
	}
}

func ListUptimeNodes() []model.UptimeNode {
	var nodes []model.UptimeNode
	DB.Find(&nodes)
	return nodes
}
