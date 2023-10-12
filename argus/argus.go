package argus

import (
	"sync"
	"time"

	"github.com/boardware-cloud/common/constants"
	"github.com/boardware-cloud/common/utils"
	argusModel "github.com/boardware-cloud/model/argus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	HEARTBEAT_TOLERANCE     = 5
	CHECK_MEMBER_INTERVAL   = 10
	HEARTBEAT_INTERVAL      = 10
	CHECK_MONITORS_INTERVAL = 10
)

type Node struct {
	Entity argusModel.ArgusNode
}

func (node *Node) Beat() {
	node.Entity.Heartbeat = time.Now().Unix()
	db.Save(&node.Entity)
}

func NewArgus(entity argusModel.Argus) Argus {
	var a Argus
	a.SetEntity(entity)
	return a
}

type Argus struct {
	entity  argusModel.Argus
	monitor Monitor
}

func (a *Argus) Alive() bool {
	entity, err := argusModel.FindArgus(a.Entity().ID)
	if err != nil || entity.Status != "ACTIVED" || entity.UpdatedAt != a.Entity().UpdatedAt {
		return false
	}
	return true
}

func (a *Argus) Spawn(node Node) {
	if !a.entity.Spawn(node.Entity.ID) {
		return
	}
	for a.Alive() {
		a.Monitor().Sleep(*a)
		if !a.Alive() {
			return
		}
		result := a.Monitor().Check()
		a.Entity().Record(string(result.Status()), result.ResponseTime())
	}
}

func (a Argus) Owner() uint {
	return a.entity.AccountId
}

func (a Argus) Monitor() Monitor {
	return a.monitor
}

func (a Argus) Entity() argusModel.Argus {
	return a.entity
}

func (a *Argus) SetEntity(m argusModel.Argus) Argus {
	a.entity = m
	a.setMonitor(m.Monitor())
	return *a
}

func (a Argus) Name() string {
	return a.entity.Name
}

func (a Argus) ID() string {
	return utils.UintToString(a.entity.ID)
}

func (a Argus) Description() string {
	return a.entity.Description
}

func (a Argus) Type() constants.MonitorType {
	return a.entity.Type
}

func (a *Argus) setMonitor(monitor argusModel.Monitor) {
	var m Monitor
	switch monitor.(type) {
	case *argusModel.HttpMonitor:
		m = &HttpMonitor{}
	case *argusModel.PingMonitor:
		m = &PingMonitor{}
	}
	m.SetEntity(monitor)
	a.monitor = m
}

type ResultStatus string

const (
	OK      ResultStatus = "OK"
	DOWN    ResultStatus = "DOWN"
	TIMEOUT ResultStatus = "TIMEOUT"
)

type Result interface {
	Status() ResultStatus
	ResponseTime() time.Duration
}

type Record struct {
	Result       ResultStatus
	ResponesTime time.Duration
	CheckedAt    time.Time
}

func Register() {
	node = new(Node)
	var mu sync.Mutex
	entity := argusModel.ArgusNode{
		Heartbeat:         time.Now().Unix(),
		HeartbeatInterval: HEARTBEAT_INTERVAL,
	}
	node.Entity = entity
	db.Save(&node.Entity)
	// Heartbeat
	go func() {
		for {
			mu.Lock()
			node.Beat()
			mu.Unlock()
			time.Sleep(HEARTBEAT_INTERVAL * time.Second)
		}
	}()
	// Recover nodes
	go func() {
		for {
			mu.Lock()
			for _, n := range diedArgusNodes() {
				go recoverNode(n)
			}
			mu.Unlock()
			time.Sleep(CHECK_MEMBER_INTERVAL * time.Second)
		}
	}()
	// Spawn argus
	go func() {
		for {
			mu.Lock()
			for _, argusEntity := range orphanArgus() {
				Spawn(NewArgus(argusEntity))
			}
			mu.Unlock()
			time.Sleep(CHECK_MONITORS_INTERVAL * time.Second)
		}
	}()
}

func orphanArgus() []argusModel.Argus {
	var argus []argusModel.Argus
	db.Find(&argus, "argus_node_id IS NULL AND status = 'ACTIVED'")
	return argus
}

func diedArgusNodes() []argusModel.ArgusNode {
	var nodes []argusModel.ArgusNode
	db.Find(&nodes)
	var diedNode []argusModel.ArgusNode
	for _, node := range nodes {
		if time.Now().Unix() > node.Heartbeat+node.HeartbeatInterval+HEARTBEAT_TOLERANCE {
			diedNode = append(diedNode, node)
		}
	}
	return diedNode
}

func recoverNode(node argusModel.ArgusNode) {
	db.Transaction(func(tx *gorm.DB) error {
		ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&argusModel.Argus{}).Where(
			"argus_node_id = ?",
			node.ID).Update("argus_node_id", nil)
		if ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		ctx = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&node)
		if ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		ctx = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Delete(&node)
		if ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		return nil
	})
}
