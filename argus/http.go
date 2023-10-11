package argus

import (
	"fmt"
	"time"

	"github.com/boardware-cloud/common/code"
	argusModel "github.com/boardware-cloud/model/argus"
)

type HttpMonitor struct {
	entity argusModel.HttpMonitor
}

func (h *HttpMonitor) Interval() time.Duration {
	return time.Duration(h.entity.Interval) * time.Second
}

func (h *HttpMonitor) Sleep() {
	time.Sleep(5 * time.Second)
	// time.Sleep(h.Interval() * time.Second)
}

func (h *HttpMonitor) SetEntity(entity argusModel.Monitor) error {
	httpMonitor, ok := entity.(*argusModel.HttpMonitor)
	if !ok {
		return code.ErrConvert
	}
	h.entity = *httpMonitor
	return nil
}

func (h *HttpMonitor) Entity() argusModel.Monitor {
	return &h.entity
}

func (h *HttpMonitor) Check() {
	fmt.Println("http check")
}

func (h *HttpMonitor) Alive() bool {
	return true
}
