package argus

import (
	"fmt"
	"sync"
	"time"

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

var db *gorm.DB

type Node struct {
	Entity argusModel.ArgusNode
}

func (node *Node) Beat() {
	node.Entity.Heartbeat = time.Now().Unix()
	db.Save(&node.Entity)
}

func (n Node) Spawn() {

}

type Argus struct {
	Monitor Monitor
}

type Monitor interface {
	Sleep()
	Interval() time.Duration
}

func Init(inject *gorm.DB) {
	db = inject
	Register(&Node{})
}

func Register(node *Node) {
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
			for _, argus := range orphanArgus() {
				go Spawn(argus)
			}
			mu.Unlock()
			time.Sleep(CHECK_MONITORS_INTERVAL * time.Second)
		}
	}()
}

func Spawn(argus argusModel.Argus) {
	var monitor Monitor
	switch argus.Monitor().GetType() {
	case "HTTP":
		m, ok := argus.Monitor().(*argusModel.HttpMonitor)
		if ok {
			monitor = HttpMonitor{
				Monitor: *m,
			}
		}
	}
	a := Argus{
		Monitor: monitor,
	}
	for {
		fmt.Println(a, "sleep")
		a.Monitor.Sleep()
	}
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
		node := argusModel.ArgusNode{}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&node, node.ID); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&argusModel.Argus{}).Where(
			"argus_node_id = ?",
			node.ID).Update("argus_node_id", nil); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		if ctx := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Delete(&argusModel.ArgusNode{}, node.ID); ctx.Error != nil {
			tx.Rollback()
			return ctx.Error
		}
		return nil
	})
}
