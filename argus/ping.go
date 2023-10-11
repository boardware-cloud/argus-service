package argus

import (
	"fmt"
	"time"

	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
)

type PingMonitor struct {
	entity argusModel.PingMonitor
}

func (p PingMonitor) Interval() time.Duration {
	return time.Duration(p.entity.Interval) * time.Second
}

func (p PingMonitor) Sleep() {
	// time.Sleep(p.Interval() * time.Second)
	time.Sleep(5 * time.Second)
}

func (p *PingMonitor) SetEntity(entity argusModel.Monitor) error {
	pingMonitor, ok := entity.(*argusModel.PingMonitor)
	if !ok {
		return code.ErrConvert
	}
	p.entity = *pingMonitor
	return nil
}

func (p PingMonitor) Entity() argusModel.Monitor {
	return &p.entity
}

func (h *PingMonitor) Check() {
	fmt.Println("ping check")
}

func (h *PingMonitor) Alive() bool {
	return false
}
